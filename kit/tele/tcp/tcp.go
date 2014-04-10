// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package tcp implements a carrier transport over TCP.
package tcp

import (
	"net"
	"strings"

	"github.com/gocircuit/circuit/kit/tele/chain"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

// Transport is a chain.Carrier over TCP.
var Transport = transport{trace.NewFrame("tcp")}

type transport struct {
	trace.Frame
}

func (transport) Listen(addr net.Addr) (net.Listener, error) {
	t := addr.String()
	if strings.Index(t, ":") < 0 {
		t = t + ":0"
	}
	l, err := net.Listen("tcp", t)
	if err != nil {
		return nil, err
	}
	return listener{l}, nil
}

func (transport) Dial(addr net.Addr) (net.Conn, error) {
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
	return &conn{trace.NewFrame("tcp", "dial"), c}, nil
}

type listener struct {
	net.Listener
}

func (l listener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return &conn{trace.NewFrame("tcp", "acpt"), c}, nil
}

type conn struct {
	trace.Frame
	net.Conn
}
