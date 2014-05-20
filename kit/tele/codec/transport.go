// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package codec

import (
	"net"
)

type CarrierTransport interface {
	Listen(addr net.Addr) CarrierListener
	Dial(addr net.Addr) CarrierConn
}

type CarrierListener interface {
	Addr() net.Addr
	Accept() CarrierConn
}

type CarrierConn interface {
	RemoteAddr() net.Addr
	Read() (chunk []byte, err error)
	Write(chunk []byte) (err error)
	Close() (err error)
}

type Transport struct {
	sub   CarrierTransport
	codec Codec
}

func NewTransport(sub CarrierTransport, codec Codec) *Transport {
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
		CarrierListener: t.sub.Listen(addr),
	}
}

type Listener struct {
	codec Codec
	CarrierListener
}

func (l *Listener) Accept() *Conn {
	return NewConn(l.CarrierListener.Accept(), l.codec)
}

func (l *Listener) Addr() net.Addr {
	return l.CarrierListener.Addr()
}
