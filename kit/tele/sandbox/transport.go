// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package sandbox provides a simulated carrier Transport for testing purposes.
package sandbox

import (
	"net"
	"sync"

	"github.com/gocircuit/circuit/kit/tele/chain"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

// Addr is a sandbox address, implementing net.Addr
type Addr string

func (a Addr) Network() string {
	return "sandbox"
}

func (a Addr) String() string {
	return string(a)
}

// Transport implements a sandbox-ed internetworking infrastructure with a customizable link behavior.
type Transport struct {
	frame trace.Frame
	connMaker
	sync.Mutex
	withAddr map[string]*listener
}

// PipeMaker is a function that creates a pipe between two addresses.
type connMaker func(af, bf trace.Frame, a, b net.Addr) (net.Conn, net.Conn)

// New creates a new sandboxed network with connections supplied by piper.
func NewTransport(f trace.Frame, connMaker connMaker) *Transport {
	return &Transport{
		frame:     f,
		connMaker: connMaker,
		withAddr:  make(map[string]*listener),
	}
}

func (s *Transport) Frame() trace.Frame {
	return s.frame
}

// Listen starts a new listener object at the given opaque address.
func (s *Transport) Listen(addr net.Addr) (net.Listener, error) {
	s.Lock()
	defer s.Unlock()
	l, ok := s.withAddr[addr.String()]
	if !ok {
		l = newListener(s.frame.Refine("listener"), addr)
		s.withAddr[addr.String()] = l
	}
	return l, nil
}

// Dial dials the opaque address.
func (s *Transport) Dial(addr net.Addr) (net.Conn, error) {
	s.Lock()
	defer s.Unlock()
	l, ok := s.withAddr[addr.String()]
	if !ok {
		return nil, chain.ErrRIP
	}
	p0, p1 := s.connMaker(
		s.frame.Refine("dial", addr.String()), s.frame.Refine("accept"),
		Addr(addr.String()), Addr(addr.String()),
	)
	l.connect(p1)
	return p0, nil
}
