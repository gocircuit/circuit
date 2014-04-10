// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package codec

import (
	"net"

	"github.com/gocircuit/circuit/kit/tele/faithful"
)

type Transport struct {
	sub   *faithful.Transport
	codec Codec
}

func NewTransport(sub *faithful.Transport, codec Codec) *Transport {
	return &Transport{sub: sub, codec: codec}
}

// Dial returns instanteneously (it does not wait on I/O operations) and always succeeds,
// returning a non-nil connection object.
func (t *Transport) Dial(addr net.Addr) *Conn {
	conn := t.sub.Dial(addr)
	return NewConn(conn, t.codec)
}

func (t *Transport) Listen(addr net.Addr) *Listener {
	return &Listener{
		codec:    t.codec,
		Listener: t.sub.Listen(addr),
	}
}

type Listener struct {
	codec Codec
	*faithful.Listener
}

func (l *Listener) Accept() *Conn {
	return NewConn(l.Listener.Accept(), l.codec)
}

func (l *Listener) Addr() net.Addr {
	return l.Listener.Addr()
}
