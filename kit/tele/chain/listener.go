// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chain

import (
	"bufio"
	"net"
	"sync"

	"github.com/gocircuit/circuit/kit/tele/limiter"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

// Listener
type Listener struct {
	frame    trace.Frame
	stat     *Stat
	listener net.Listener
	withID__ sync.Mutex
	withID   map[chainID]*acceptConn
	accpt__  sync.Mutex
	accpt    chan *Conn
}

const (
	MaxHandshakes = 5 // Maximum number of concurrent handshake interactions
)

func NewListener(frame trace.Frame, stat *Stat, carrier Carrier, addr net.Addr) *Listener {
	cl, err := carrier.Listen(addr)
	if err != nil {
		panic(err)
	}
	l := &Listener{
		frame:    frame,
		stat:     stat,
		listener: cl,
		withID:   make(map[chainID]*acceptConn),
		accpt:    make(chan *Conn),
	}
	l.frame.Bind(l)
	go l.loop()
	return l
}

func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}

func (l *Listener) loop() {
	lmtr := limiter.New(MaxHandshakes)
	for {
		lmtr.Open()
		c, err := l.listener.Accept()
		if err != nil {
			lmtr.Close()
			panic(err) // Best not to be quiet about it
		}
		go func() {
			l.handshake(c)
			lmtr.Close()
		}()
	}
}

func (l *Listener) handshake(carrier net.Conn) {
	r := bufio.NewReader(carrier)
	dialMsg, err := readMsgDial(r)
	if err != nil {
		l.frame.Printf("handshake unrecognized dial message (%s)", err)
		makeWelcome(rejectDial).Write(carrier) // reject chain permanently
		carrier.Close()
		return
	}

	switch dialMsg.SeqNo {
	case 0:
		l.frame.Printf("rejecting connection with sequence number 0 (remote bug)")
		makeWelcome(rejectZero).Write(carrier) // reject chain permanently
		carrier.Close()
		return

	case 1:
		var ac *acceptConn
		if ac, err = l.make(dialMsg.ID, carrier, r); err != nil {
			l.frame.Printf("rejecting duplicate connection (%s) %x", err, dialMsg.ID)
			makeWelcome(rejectDup).Write(carrier) // reject chain permanently
			carrier.Close()
			return
		}
		if err = makeWelcome(rejectOK).Write(carrier); err != nil {
			l.frame.Printf("carrier broke during first welcome (%s), ChainID=%x", err, dialMsg.ID)
			carrier.Close()
			return
		}
		// Send connection for Accept, if not a re-dial
		l.accpt__.Lock()
		l.accpt <- &ac.Conn
		l.accpt__.Unlock()

	default:
		var ac *acceptConn
		if ac = l.get(dialMsg.ID); ac == nil {
			l.frame.Printf("rejecting redial on closed chain, ChainID=%x", dialMsg.ID)
			// This is a valid condition, not protocol misbehavior, so we need to notify the
			// remote that they are talking to a new process with a different identity.
			makeWelcome(rejectClosed).Write(carrier) // reject chain permanently
			carrier.Close()
			return
		}
		if err = makeWelcome(rejectOK).Write(carrier); err != nil {
			l.frame.Printf("carrier broke during follow-on welcome (%s), ChainID=%x", err, dialMsg.ID)
			carrier.Close()
			return
		}
		ac.Accept(carrier, r, dialMsg.SeqNo)
	}
}

func (l *Listener) Accept() *Conn {
	return <-l.accpt
}

func (l *Listener) make(id chainID, carrier net.Conn, r *bufio.Reader) (*acceptConn, error) {
	l.withID__.Lock()
	defer l.withID__.Unlock()
	if _, ok := l.withID[id]; ok {
		return nil, errDup
	}
	addr := carrier.RemoteAddr()
	l.stat.addAC(id, addr)
	ac := newAcceptConn(l.frame.Refine(l.Addr().String()), id, addr, carrier, r, func() {
		l.scrub(id)
		l.stat.scrubAC(id)
	})
	l.withID[id] = ac
	return ac, nil
}

func (l *Listener) get(id chainID) *acceptConn {
	l.withID__.Lock()
	defer l.withID__.Unlock()
	return l.withID[id]
}

func (l *Listener) scrub(id chainID) *acceptConn {
	l.withID__.Lock()
	defer l.withID__.Unlock()
	a, ok := l.withID[id]
	if !ok {
		return nil
	}
	delete(l.withID, id)
	return a
}
