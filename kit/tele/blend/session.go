// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package blend

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/gocircuit/circuit/kit/tele/codec"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

// AcceptSession
type AcceptSession struct {
	Session
}

func newAcceptSession(frame trace.Frame, sub *codec.Conn) *AcceptSession {
	ab := &AcceptSession{}
	ab.init(frame, make(chan *Conn, 1), false, sub, nil)
	return ab
}

func (as *AcceptSession) Accept() *Conn {
	return <-as.Session.ach
}

// DialSession
type DialSession struct {
	Session
}

func newDialSession(frame trace.Frame, sub *codec.Conn, scrb func()) *DialSession {
	db := &DialSession{}
	db.init(frame, nil, true, sub, scrb)
	return db
}

func (ds *DialSession) Dial() *Conn {
	ds.Session.o_.Lock()
	defer ds.Session.o_.Unlock()
	if ds.Session.o_open[ds.o_ndial] != nil {
		panic("u")
	}
	conn := newConn(ds.Session.o_ndial, &ds.Session)
	ds.Session.o_open[ds.Session.o_ndial] = conn
	ds.Session.o_ndial++
	return conn
}

// Session
type Session struct {
	frame   trace.Frame
	scrb    func()
	ach     chan *Conn
	dialing bool
	subr    *codec.Conn
	//
	o_      sync.Mutex // Sync o_open conns structure
	o_ndial ConnID
	o_open  map[ConnID]*Conn
	o_use   time.Time
	//
	w_    sync.Mutex // Linearize write ops on sub
	w_sub *codec.Conn
}

func (ssn *Session) init(frame trace.Frame, accept chan *Conn, dialing bool, sub *codec.Conn, scrb func()) {
	ssn.frame = frame
	ssn.scrb = scrb
	ssn.frame.Bind(ssn)
	ssn.ach, ssn.dialing = accept, dialing
	ssn.subr, ssn.w_sub = sub, sub
	ssn.o_open = make(map[ConnID]*Conn)
	ssn.o_use = time.Now()

	go ssn.readloop()
}

func (ssn *Session) String() string {
	ssn.w_.Lock()
	defer ssn.w_.Unlock()
	if ssn.w_sub == nil {
		return "session closed"
	}
	return ssn.w_sub.String()
}

func (ssn *Session) NumConn() (numconn int, lastuse time.Time) {
	ssn.o_.Lock()
	defer ssn.o_.Unlock()
	return len(ssn.o_open), ssn.o_use
}

func (ssn *Session) RemoteAddr() net.Addr {
	return ssn.subr.RemoteAddr()
}

func (ssn *Session) hijack() (w_sub *codec.Conn) {
	ssn.w_.Lock()
	defer ssn.w_.Unlock()
	w_sub, ssn.w_sub = ssn.w_sub, nil
	return
}

func (ssn *Session) Close() (err error) {
	w_sub := ssn.hijack()
	if w_sub == nil {
		return io.ErrClosedPipe
	}
	return w_sub.Close()
}

func (ssn *Session) teardown() {
	// Notify accepters, if an accept session
	if ssn.ach != nil {
		close(ssn.ach)
	}

	ssn.o_.Lock()
	// The substrate connection does not allow Write after Close.
	// To prevent writes from Conns hitting the substrate before the Conns have been notified:
	// we first remove the substrate from its field to prevents writes from Conn going through to it,
	// and then we close the substrate.
	if w_sub := ssn.hijack(); w_sub != nil {
		w_sub.Close()
	}
	// Notify o_open connections
	for connID, conn := range ssn.o_open {
		conn.promptClose()
		delete(ssn.o_open, connID)
	}
	ssn.o_.Unlock()
	if ssn.scrb != nil {
		ssn.scrb()
	}
}

func (ssn *Session) readloop() {
	defer ssn.teardown()
	for {
		if err := ssn.read(); err != nil {
			// ssn.frame.Printf("session read loop (%s)", err)
			return
		}
	}
}

func (ssn *Session) read() error {
	msg := &Msg{}
	if err := ssn.subr.Read(msg); err != nil {
		// Connection broken
		return err
	}

	switch t := msg.Demux.(type) {
	case *PayloadMsg:
		conn := ssn.get(msg.ConnID)
		if conn != nil {
			// Existing connection
			conn.prompt(t.Payload, nil)
			return nil
		}
		// Dead connection
		if t.SeqNo > 0 {
			ssn.writeAbort(msg.ConnID, ErrGone)
			return nil
		}
		// New connection
		if ssn.ach != nil {
			conn = newConn(msg.ConnID, ssn)
			ssn.set(msg.ConnID, conn)
			conn.prompt(t.Payload, nil)
			ssn.ach <- conn // Send new connection to user
			return nil
		} else {
			ssn.writeAbort(msg.ConnID, ErrOff)
			return nil
		}

	case *AbortMsg:
		conn := ssn.get(msg.ConnID)
		if conn == nil {
			// Discard CLOSE for non-existent connections
			// Do not respond with a CLOSE packet. It would cause an avalanche of CLOSEs.
			return nil
		}
		ssn.scrub(msg.ConnID)
		conn.prompt(nil, t.Err)
		return nil
	}

	// Unexpected remote behavior
	return ErrClash
}

func (ssn *Session) count() int {
	ssn.o_.Lock()
	defer ssn.o_.Unlock()
	return len(ssn.o_open)
}

func (ssn *Session) get(connID ConnID) *Conn {
	ssn.o_.Lock()
	defer ssn.o_.Unlock()
	ssn.o_use = time.Now()
	return ssn.o_open[connID]
}

func (ssn *Session) set(connID ConnID, conn *Conn) {
	ssn.o_.Lock()
	defer ssn.o_.Unlock()
	ssn.o_open[connID] = conn
}

func (ssn *Session) scrub(connID ConnID) {
	ssn.o_.Lock()
	defer ssn.o_.Unlock()
	delete(ssn.o_open, connID)
}

func (ssn *Session) write(msg *Msg) error {
	ssn.w_.Lock()
	defer ssn.w_.Unlock()
	//
	if ssn.w_sub == nil {
		return io.ErrUnexpectedEOF
	}
	if err := ssn.w_sub.Write(msg); err != nil {
		return err
	}
	return nil
}

func (ssn *Session) writeAbort(connID ConnID, reason error) error {
	msg := &Msg{
		ConnID: connID,
		Demux: &AbortMsg{
			Err: reason,
		},
	}
	return ssn.write(msg)
}
