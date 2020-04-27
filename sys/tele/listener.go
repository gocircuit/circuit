// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tele

import (
	"encoding/gob"
	"log"
	"net"
	"os"
	"sync"

	"github.com/hoijui/circuit/kit/tele/blend"
	"github.com/hoijui/circuit/use/errors"
	"github.com/hoijui/circuit/use/n"
)

// Listener
type Listener struct {
	addr     *Addr
	listener *blend.Listener
	ach__    sync.Mutex
	ach      chan n.Conn
}

func newListener(workerID n.WorkerID, pid int, listener *blend.Listener) *Listener {
	l := &Listener{
		addr:     NewNetAddr(workerID, pid, listener.Addr()), // Compute what our address looks like on the outside.
		listener: listener,
		ach:      make(chan n.Conn),
	}
	go l.loop()
	return l
}

func (l *Listener) loop() {
	for {
		session := l.listener.AcceptSession()
		go func() {
			defer session.Close()
			// Authenticate dialer on first connection
			sourceAddr, err := l.handshake(session.Accept())
			if err != nil {
				return
			}
			for {
				conn := session.Accept()
				if conn == nil {
					// Nil conn signifies the session has been closed
					return
				}
				// For now, listening sessions do not expire themselves on inactivity to prevent
				// race against DialSessions, who currently hold the sole responsibility.
				l.ach__.Lock()
				l.ach <- NewConn(conn, sourceAddr)
				l.ach__.Unlock()
			}
		}()
	}
}

func (l *Listener) handshake(conn *blend.Conn) (sourceAddr *Addr, err error) {
	if conn == nil {
		return nil, errors.NewError("listener off")
	}
	defer conn.Close()
	//
	var msg interface{}
	msg, err = conn.Read()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			conn.Write(&RejectMsg{err})
		} else {
			err = conn.Write(&WelcomeMsg{})
		}
	}()
	hello, ok := msg.(*HelloMsg)
	if !ok {
		log.Println("rejecting", conn.RemoteAddr().String(), "unknown hello message type")
		return nil, errors.NewError("rejecting unknown hello type")
	}
	// Accept user connections
	da, ok := hello.SourceAddr.(*Addr)
	if !ok {
		log.Println("rejecting", conn.RemoteAddr().String(), "unknown source address type")
		return nil, errors.NewError("rejecting unknown source address type")
	}
	reverseAddr(da, conn.RemoteAddr())
	la, ok := hello.TargetAddr.(*Addr)
	if !ok {
		log.Println("rejecting ", conn.RemoteAddr().String(), "unknown target address type")
		return nil, errors.NewError("rejecting unknown target address type")
	}
	if la.WorkerID() != l.addr.WorkerID() {
		log.Println("rejecting", conn.RemoteAddr().String(), "due to worker identity mismatch")
		return nil, errors.NewError("rejecting worker identity mismatch, looks for %s, got %s", la.WorkerID(), l.addr.WorkerID())
	}
	if la.PID != os.Getpid() {
		log.Println("rejecting", conn.RemoteAddr().String(), "due to worker PID mismatch")
		return nil, errors.NewError("rejecting worker PID mismatch, looks for %d, got %d", la.PID, os.Getpid())
	}
	return da, nil
}

func reverseAddr(bound *Addr, seen net.Addr) {
	// var saved = bound.String()
	if !bound.TCP.IP.IsUnspecified() {
		return
	}
	bound.TCP.IP = seen.(*net.TCPAddr).IP
	// log.Printf("Reverse dial address auto-completed: %s => %s", saved, bound.String())
}

func (l *Listener) Accept() n.Conn {
	return <-l.ach
}

func (l *Listener) Addr() n.Addr {
	return l.addr
}

func (l *Listener) Close() {}

// Dialer sends HelloMsg to accepter when opening a session to advertise its workerID and local process ID
type HelloMsg struct {
	SourceAddr n.Addr
	TargetAddr n.Addr
}

type WelcomeMsg struct{}

type RejectMsg struct {
	Err error
}

func init() {
	gob.Register(&HelloMsg{})
	gob.Register(&WelcomeMsg{})
	gob.Register(&RejectMsg{})
}
