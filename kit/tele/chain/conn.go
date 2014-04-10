// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chain

import (
	"bufio"
	"io"
	"net"
	"sync"

	"github.com/gocircuit/circuit/kit/tele/trace"
)

// Conn is a chunk connection that transparently repairs connectivity if the underlying net.Conn breaks.
type Conn struct {
	frame    trace.Frame
	scrb     func()
	id       chainID
	addr     net.Addr
	linker   linker
	cascade  *Cascade         // Keeps track of current ConnWriter
	rch      chan interface{} // readLoop -> Read
	k__      sync.Mutex
	killed   bool
	kch      chan struct{} // Kill channel
	kcarrier net.Conn      // Pointer to current carrier connection
}

type delivered interface {
	Delivered() bool
}

type linker interface {
	// If error is non-nil, the destination is permanently gone.
	Link(error) (net.Conn, *bufio.Reader, SeqNo, error)
	Kill()
}

func (a *Conn) Start(frame trace.Frame, id chainID, addr net.Addr, linker linker, scrb func()) {
	frame.Bind(a)
	a.frame = frame
	a.scrb = scrb
	a.id = id
	a.addr = addr
	a.linker = linker
	a.cascade = MakeCascade(frame)
	// A buffer size 1 on rch, helps remove a deadlock in the TestConn.
	// Essentially it ensures that Read and Write (on two ends of a
	// connection) cannot deadlock each other when a successful Write also
	// requires a stitch. We throw in a couple of extra buffer spaces to
	// prevent any potential deadlock between Read and Kill.
	a.rch = make(chan interface{}, 3)
	a.kch = make(chan struct{})
	go a.readLoop()
}

// RemoteAddr returns the address of the remote endpoint of this connection.
func (a *Conn) RemoteAddr() net.Addr {
	return a.addr
}

type debugger interface {
	Debug() interface{}
}

func (a *Conn) Debug() interface{} {
	iv := a.cascade.Current()
	if iv == nil {
		return nil
	}
	cw, ok := iv.Value().(*ConnWriter)
	if !ok {
		return nil
	}
	udbg, ok := cw.carrier.(debugger)
	if !ok {
		return nil
	}
	return udbg.Debug()
}

// Read reads the next chunk of bytes.
// Return table:
//
//		[]byte		error
//		------      -----
//		chunk		nil						=> Chunk received on same carrier connection as last Read
//		nil			ErrStitch				=> Carrier writer connection changed, no new messages received
//		nil			non-nil, non-ErrStitch	=> Connection permanently terminated
//
// Read must be called greedily, even if no payload chunks are expected from
// the other side, as it may have to report carrier connection stitching events.
func (a *Conn) Read() (v []byte, err error) {
	msg, ok := <-a.rch
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}
	if err, ok = msg.(error); ok {
		return nil, err
	}
	return msg.([]byte), nil
}

// The dialer endpoint of a connection must be able to detect connection outage in order to redial.
// Canonically, communication errors are detected during Read operations. (This cannot always be done reliably on write.)
// Thus, the dialer connection must spawn a dedicate a read loop, in case the user logic does not use Read.
func (a *Conn) readLoop() {
	var (
		carrier net.Conn
		r       *bufio.Reader
		seqno   SeqNo
		err     error
		v       []byte
	)
	defer func() {
		close(a.rch)
		a.cascade.Close()
		a.linker.Kill()
		if carrier != nil {
			carrier.Close()
		}
		if a.scrb != nil {
			a.scrb()
		}
	}()
	// Loop
	for {
		// Check for kill signal
		select {
		case <-a.kch:
			return
		default:
		}
		// Read from connection
		v, err = a.read(r)
		if err == nil {
			a.rch <- v
			continue
		}
		// Link
		if carrier != nil {
			carrier.Close()
			r = nil
		}
		if carrier, r, seqno, err = a.linker.Link(err); err != nil {
			// Permanently cannot reconnect
			return
		}
		a.setKillCarrier(carrier)
		w := &ConnWriter{conn: a, carrier: carrier}
		a.cascade.Transition(w)
		a.rch <- &ErrStitch{SeqNo: seqno, Writer: w}
	}
}

func (a *Conn) read(r *bufio.Reader) (v []byte, err error) {
	if r == nil {
		return nil, io.ErrUnexpectedEOF
	}
	msg, err := readMsgPayload(r)
	if err != nil {
		if msg != nil {
			panic("eh")
		}
		return nil, err
	}
	if len(msg.Payload) == 0 {
		panic("eh")
	}
	return msg.Payload, nil
}

type ConnWriter struct {
	conn      *Conn
	carrier   net.Conn
	w__       sync.Mutex
	delivered bool // Was the last write successful; only with sandbox underlying
}

// Write writes the chunk v to the connection.
// If Write returns a non-nil error, the underlying writer connection is broken permanently.
// A new one should be obtained from Conn.Read.
func (a *ConnWriter) Write(v []byte) (err error) {
	a.w__.Lock()
	defer a.w__.Unlock()
	err = (&msgPayload{v}).Write(a.carrier)
	if oob, ok := a.carrier.(delivered); ok {
		a.delivered = oob.Delivered()
	}
	return err
}

// Delivered returns true if the last Write was successfully delivered to the destination.
// Delivered is supported only when the carrier transport is a sandbox, and should be used
// only for test purposes.
func (a *ConnWriter) Delivered() bool {
	a.w__.Lock()
	defer a.w__.Unlock()
	return a.delivered
}

// Kill closes the connection permanently and not gracefully. Chain connections cannot be closed gracefully.
func (a *Conn) Kill() {
	a.k__.Lock()
	defer a.k__.Unlock()
	if a.killed {
		return
	}
	if a.kcarrier != nil {
		a.kcarrier.Close()
	}
	a.linker.Kill() // Unblock readLoop if blocked on linker.Link
	close(a.kch)
	a.killed = true
}

func (a *Conn) setKillCarrier(kcarrier net.Conn) {
	a.k__.Lock()
	defer a.k__.Unlock()
	a.kcarrier = kcarrier
}

func (a *Conn) Close() error {
	panic("chain does not support graceful close")
}
