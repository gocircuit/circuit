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
	"time"

	"github.com/gocircuit/circuit/kit/tele/trace"
)

/*
	Kill (turns Accept into a nop)
	  |
	  V
	Accept -- ch --> Link

*/

type acceptConn struct {
	Conn
	ach      chan *accept
	accept__ sync.Mutex
	closed   bool
	link__   sync.Mutex
	seqno    SeqNo // Sequence number of the current underlying connection
	carrier  net.Conn
}

type accept struct {
	Carrier net.Conn
	R       *bufio.Reader
	SeqNo   SeqNo
}

func newAcceptConn(frame trace.Frame, id chainID, addr net.Addr, carrier net.Conn, r *bufio.Reader, scrb func()) *acceptConn {
	ac := &acceptConn{
		ach: make(chan *accept, MaxHandshakes+3),
	}
	ac.Conn.Start(frame, id, addr, (*acceptLink)(ac), scrb)
	ac.Accept(carrier, r, 1)
	return ac
}

func (ac *acceptConn) Accept(carrier net.Conn, r *bufio.Reader, seqno SeqNo) {
	ac.accept__.Lock()
	defer ac.accept__.Unlock()
	if ac.closed {
		carrier.Close()
		return
	}
	ac.ach <- &accept{carrier, r, seqno}
}

type acceptLink acceptConn

// Kill shuts down the acceptLink, interrupting a pending wait for connection in Link.
func (al *acceptLink) Kill() {
	al.accept__.Lock()
	defer al.accept__.Unlock()
	if al.closed {
		return
	}
	al.closed = true
	close(al.ach)
}

// Link blocks until a new connection to the remote endpoint is passed through Accept or Kill is invoked.
func (al *acceptLink) Link(reason error) (net.Conn, *bufio.Reader, SeqNo, error) {
	al.link__.Lock()
	defer al.link__.Unlock()
	if al.carrier != nil {
		al.carrier.Close()
	}
	for {
		select {
		case <-time.After(7*time.Second): 
			// If the dialer does not attempt to recover the connection within 7 seconds,
			// the listener (us) drops it.
			// TODO: If the reason for the disconnect is a temporary network partition,
			// the dialer won't be able to recover this particular chain, but higher-level
			// logic should retry connecting  to the same address on a new connection once,
			// since this should succeed as a new chain connection after the partition is over.
			return nil, nil, 0, io.ErrUnexpectedEOF
		case replaceWith, ok := <-al.ach:
			if !ok {
				return nil, nil, 0, io.ErrUnexpectedEOF
			}
			if replaceWith.SeqNo > al.seqno {
				al.seqno = replaceWith.SeqNo
				al.carrier = replaceWith.Carrier
				//al.Conn.frame.Printf("ACCEPTED #%d", al.seqno)
				return al.carrier, replaceWith.R, al.seqno, nil
			}
			al.Conn.frame.Printf("out-of-order redial #%d arrived while using #%d", replaceWith.SeqNo, al.seqno)
			replaceWith.Carrier.Close()
		}
	}
	panic("u")
}
