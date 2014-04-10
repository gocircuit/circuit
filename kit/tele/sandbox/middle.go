// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package sandbox

import (
	"net"
	"sync"
	"time"

	"github.com/gocircuit/circuit/kit/tele/trace"
)

// writeCloseRace is the buffer size of the sandbox message channel, chosen so that
// it would absorb a few Writes that are competing with a Close on another thread.
const writeCloseRace = 3

func newChan() chan interface{} {
	return make(chan interface{}, writeCloseRace)
}

type feedback struct {
	OK  bool
	EOF bool
}

func newBackChan() chan feedback {
	return make(chan feedback)
}

// NewPipe returns the two net.Conn endpoints of a reliable pipe.
func NewPipe(f0, f1 trace.Frame, a0, a1 net.Addr) (p, q net.Conn) {
	a, b := NewHalfConn(f0), NewHalfConn(f1)
	ab, ba := newChan(), newChan()
	a.RecvFrom(ba, a1)
	b.RecvFrom(ab, a0)
	a.SendTo(ab, nil, a0)
	b.SendTo(ba, nil, a1)
	return a, b
}

// NewSievePipe returns the two net.Conn endpoints of a new bi-directional drop tail pipe.
func NewSievePipe(f0, f1 trace.Frame, a0, a1 net.Addr, nok, ndrop int, expa, expb time.Duration) (p, q net.Conn) {
	a, b := NewHalfConn(f0), NewHalfConn(f1)
	ax, bx, xa, xb := newChan(), newChan(), newChan(), newChan()
	axk, bxk := newBackChan(), newBackChan()
	a.RecvFrom(xa, a1)
	b.RecvFrom(xb, a0)
	a.SendTo(ax, axk, a0)
	b.SendTo(bx, bxk, a1)
	ab, ba := StartSieve(f0.Refine("sieve"), f1.Refine("sieve"), ax, bx, xa, xb, axk, bxk, nok, ndrop, expa, expb)
	a.in, a.out = ba, ab
	b.in, b.out = ab, ba
	return a, b
}

// bi converts bool to int
func bi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Send-closer protects send and close operations from concurring agents to a channel.
type sendCloser struct {
	sync.Mutex
	Ch chan<- interface{}
}

func newSendCloser(send chan<- interface{}) *sendCloser {
	return &sendCloser{Ch: send}
}

func (x *sendCloser) Init(send chan<- interface{}) {
	x.Ch = send
}

func (x *sendCloser) Send(p interface{}) bool {
	x.Lock()
	defer x.Unlock()
	if x.Ch == nil {
		return false
	}
	x.Ch <- p
	return true
}

func (x *sendCloser) Close() {
	x.Lock()
	defer x.Unlock()
	if x.Ch == nil {
		return
	}
	close(x.Ch)
	x.Ch = nil
}

/*
	+---+          +--------------+          +---+
	|   |<---xa----t    Sieve     s<---bx----|   |
	|   |          |              k----bxk-->|   |
	|   |          +---|C|---+----+          |   |
	|   |               ^    |               |   |
	| A |               |    |               | B |
	|   |               |    V               |   |
	|   |          +----+---|C|---|          |   |
	|   |<---axk---k              |          |   |
	|   |----ax--->s    Sieve     t----xb--->|   |
	+---+          +--------------+          +---+
	HalfConn                                 HalfConn
*/

func StartSieve(
	fab, fba trace.Frame,
	ax, bx <-chan interface{},
	xa, xb chan<- interface{},
	axk, bxk chan<- feedback,
	nok, ndrop int,
	expa, expb time.Duration,
) (ab, ba *Sieve) {
	ab, ba = &Sieve{
		Frame: fab,
		abrt:  make(chan struct{}),
		exp:   expa,
		nok:   nok,
		ndrop: ndrop,
		s:     ax,
		k:     axk,
	}, &Sieve{
		Frame: fba,
		abrt:  make(chan struct{}),
		exp:   expb,
		nok:   nok,
		ndrop: ndrop,
		s:     bx,
		k:     bxk,
	}
	ab.Frame.Bind(ab)
	ba.Frame.Bind(ba)
	ab.t.Init(xb)
	ba.t.Init(xa)
	ab.dual, ba.dual = ab, ba
	go ab.loop()
	go ba.loop()
	return
}

/*
	       +---+
	<--k---|   |
	       |   +---+
	---s-->| Sieve |---t-->
	       +-------+
*/
type Sieve struct {
	trace.Frame
	abrt  chan struct{}
	exp   time.Duration      // Sieve connctions are closed after approximately exp duration of inactivity
	nok   int                // Number of writes to transmit faithfully
	ndrop int                // Number of writes to drop right before breaking the connection
	s     <-chan interface{} // Incoming writes
	k     chan<- feedback    // Feedback to writer
	t     sendCloser         // Outgoing reads
	dual  *Sieve             // The other side of the bi-directional connection
}

// TIMEOUT MECHANISM:
// If the last packet was dropped, no eof has been reached and some time has
// passed since the last write to the connection, then kill the connection. (A
// user-side successful pre-EOF write, which was actually dropped, might cause
// a deadlock if it was the last write, as the recipient will block on read
// forever.) For this reason, we have a timeout mechanism inside the Sieve,
// which closes the connection after a period of inactivity.

func (x *Sieve) loop() {
	// Close out connection
	defer func() {
		close(x.k)
		x.Mute()
		x.dual.Mute()
	}()
	// Timeout ticker
	var tkr *time.Ticker
	var expchan <-chan time.Time
	if x.exp > 0 {
		tkr = time.NewTicker(x.exp)
		expchan = tkr.C
	}
	defer func() {
		if tkr != nil {
			tkr.Stop()
		}
	}()
	// Main loop
	var nrecv, nticks int
	for {
		select {
		case <-x.abrt:
			x.Frame.Println("Aborting sieve.")
			return
		case <-expchan:
			nticks++
			if nticks > 1 {
				x.Frame.Println("Breaking connection due to timeout.")
				return // Timeout with nothing send kills the connection
			}
		case q, ok := <-x.s:
			nticks = 0
			if !ok {
				return // s-channel closed, i.e. a Close on the sending HalfConn, kills the connection
			}
			nrecv++
			send := nrecv <= x.nok
			forceEOF := nrecv+1 > x.nok+x.ndrop
			if send {
				if _, ok := q.(eof); ok {
					x.Frame.Println("SIEVE ---> EOF") // ??
				}
				x.t.Send(q)
			}
			// OK tells conn.Write whether the message was delivered
			// EOF tells conn.Write whether the connection has ended after the write
			x.k <- feedback{OK: send, EOF: forceEOF}
			if forceEOF {
				x.Mute()
				x.dual.Mute()
			}
			//x.Printf("X=%-2d dlvr=%4v forceEOF=%4v (%d,%d)", nrecv, send, forceEOF, x.nok, x.ndrop)
		}
	}
}

// Mute closes the sink-side of the pipe, alerting the read-side of the sink that the connection is broken.
func (x *Sieve) Mute() {
	x.t.Close()
}

// Abort can be invoked only once.
func (x *Sieve) Abort() {
	close(x.abrt)
}
