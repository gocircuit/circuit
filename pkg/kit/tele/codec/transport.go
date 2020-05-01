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
	Dial(addr net.Addr) (CarrierConn, error)
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
func (t *Transport) Dial(addr net.Addr) (*Conn, error) {
	conn, err := t.sub.Dial(addr)
	if err != nil {
		return nil, err
	}
	return NewConn(conn, t.codec), nil
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
	for {
		conn := l.CarrierListener.Accept()
		if conn == nil {
			continue
		}
		return NewConn(conn, l.codec)
	}
	panic(0)
}

func (l *Listener) Addr() net.Addr {
	return l.CarrierListener.Addr()
}
