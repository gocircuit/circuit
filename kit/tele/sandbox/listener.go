// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package sandbox

import (
	"io"
	"net"
	"sync"

	"github.com/gocircuit/circuit/kit/tele/trace"
)

// listener implements a net.Listener for the sandbox Transport.
type listener struct {
	trace.Frame
	addr net.Addr
	ch__ sync.Mutex
	ch   chan net.Conn
}

func newListener(f trace.Frame, addr net.Addr) *listener {
	l := &listener{Frame: f, addr: addr, ch: make(chan net.Conn)}
	l.Frame.Bind(l)
	return l
}

func (sl *listener) connect(p net.Conn) {
	sl.ch__.Lock()
	defer sl.ch__.Unlock()
	sl.ch <- p
}

func (sl *listener) Accept() (net.Conn, error) {
	p, ok := <-sl.ch
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}
	return p, nil
}

func (sl *listener) Close() error {
	sl.ch__.Lock()
	defer sl.ch__.Unlock()
	close(sl.ch)
	sl.ch = nil
	return nil
}

func (sl *listener) Addr() net.Addr {
	return sl.addr
}
