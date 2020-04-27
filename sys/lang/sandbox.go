// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"encoding/gob"
	"io"
	"net"
	"sync"

	"github.com/hoijui/circuit/use/n"
)

type sandbox struct {
	lk sync.Mutex
	l  map[n.WorkerID]*listener
}

var s = &sandbox{l: make(map[n.WorkerID]*listener)}

// NewSandbox creates a new transport instance, part of a sandbox network in memory
func NewSandbox() n.Transport {
	s.lk.Lock()
	defer s.lk.Unlock()

	l := &listener{
		id: n.ChooseWorkerID(),
		ch: make(chan *halfconn),
	}
	l.a = &addr{ID: l.id, l: l}
	s.l[l.id] = l
	return l
}

func (l *listener) Listen(net.Addr) n.Listener {
	return l
}

func dial(remote n.Addr) (n.Conn, error) {
	pr, pw := io.Pipe()
	qr, qw := io.Pipe()
	srvhalf := &halfconn{PipeWriter: qw, PipeReader: pr}
	clihalf := &halfconn{PipeWriter: pw, PipeReader: qr}
	s.lk.Lock()
	l := s.l[remote.(*addr).WorkerID()]
	s.lk.Unlock()
	if l == nil {
		panic("unknown listener id")
	}
	go func() {
		l.ch <- srvhalf
	}()
	return ReadWriterConn(l.Addr(), clihalf), nil
}

// addr implements Addr
type addr struct {
	ID n.WorkerID
	l  *listener
}

func (a *addr) WorkerID() n.WorkerID {
	return a.ID
}

func (a *addr) NetAddr() net.Addr {
	return a
}

func (a *addr) Network() string {
	return "sandbox"
}

func (a *addr) String() string {
	return a.ID.String()
}

func (a *addr) FileName() string {
	return a.ID.String()
}

func init() {
	gob.Register(&addr{})
}

// listener implements Listener
type listener struct {
	id n.WorkerID
	a  *addr
	ch chan *halfconn
}

func (l *listener) Addr() n.Addr {
	return l.a
}

func (l *listener) Accept() n.Conn {
	return ReadWriterConn(l.Addr(), <-l.ch)
}

func (l *listener) Close() {
	s.lk.Lock()
	defer s.lk.Unlock()
	delete(s.l, l.id)
}

func (l *listener) Dial(remote n.Addr) (n.Conn, error) {
	return dial(remote)
}

// halfconn is one end of a byte-level connection
type halfconn struct {
	*io.PipeReader
	*io.PipeWriter
}

func (h *halfconn) Close() error {
	return h.PipeWriter.Close()
}
