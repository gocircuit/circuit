// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package tcp implements carrier transports over TCP.
package tcp

import (
	"net"
	"strings"

	"github.com/gocircuit/circuit/kit/tele/chain"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

// ChainTransport is a chain.Carrier over TCP.
var ChainTransport = chainTransport{trace.NewFrame("tcp")}

type chainTransport struct {
	trace.Frame
}

func (chainTransport) Listen(addr net.Addr) (net.Listener, error) {
	t := addr.String()
	if strings.Index(t, ":") < 0 {
		t = t + ":0"
	}
	l, err := net.Listen("tcp", t)
	if err != nil {
		return nil, err
	}
	return chainListener{l}, nil
}

func (chainTransport) Dial(addr net.Addr) (net.Conn, error) {
	c, err := net.Dial("tcp", addr.String())
	if err != nil {
		operr, ok := err.(*net.OpError)
		if !ok {
			return nil, err
		}
		if operr.Temporary() {
			return nil, err
		}
		return nil, chain.ErrRIP
	}
	return newChainConn(trace.NewFrame("tcp", "dial"), c.(*net.TCPConn)), nil
}

type chainListener struct {
	net.Listener
}

func (l chainListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return newChainConn(trace.NewFrame("tcp", "acpt"), c.(*net.TCPConn)), nil
}

type chainConn struct {
	trace.Frame
	*net.TCPConn
}

func newChainConn(f trace.Frame, c *net.TCPConn) *chainConn {
	if err := c.SetKeepAlive(true); err != nil {
		panic(err)
	}
	return &chainConn{f, c}
}
