// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package faithful

import (
	"net"

	"github.com/gocircuit/circuit/kit/tele/chain"
	"github.com/gocircuit/circuit/kit/tele/trace"
	"github.com/gocircuit/circuit/kit/tele/codec"
)

// Transport
type Transport struct {
	frame trace.Frame
	chain *chain.Transport
}

func NewTransport(frame trace.Frame, chain *chain.Transport) codec.CarrierTransport {
	t := &Transport{frame: frame, chain: chain}
	t.frame.Bind(t)
	return t
}

func (t *Transport) Listen(addr net.Addr) codec.CarrierListener {
	return NewListener(t.frame.Refine("listener"), t.chain.Listen(addr))
}

// Dial returns instanteneously (it does not wait on I/O operations) and always succeeds,
// returning a non-nil connection object.
func (t *Transport) Dial(addr net.Addr) codec.CarrierConn {
	conn := t.chain.Dial(addr)
	return NewConn(t.frame.Refine("dial"), conn)
}

// Listener
type Listener struct {
	frame trace.Frame
	sub   *chain.Listener
}

func NewListener(f trace.Frame, sub *chain.Listener) codec.CarrierListener {
	l := &Listener{frame: f, sub: sub}
	l.frame.Bind(l)
	return l
}

func (l *Listener) Addr() net.Addr {
	return l.sub.Addr()
}

func (l *Listener) Accept() codec.CarrierConn {
	return NewConn(l.frame.Refine("accept"), l.sub.Accept())
}
