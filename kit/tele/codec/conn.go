// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package codec

import (
	"net"
	"io"
)

type Conn struct {
	enc   Encoder
	dec   Decoder
	carrier CarrierConn
}

func NewConn(carrier CarrierConn, codec Codec) *Conn {
	if carrier == nil {
		panic("nil carrier")
	}
	return &Conn{
		enc:   codec.NewEncoder(),
		dec:   codec.NewDecoder(),
		carrier: carrier,
	}
}

func (c *Conn) String() string {
	return c.RemoteAddr().String()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.carrier.RemoteAddr()
}

func (c *Conn) Write(v interface{}) (err error) {
	chunk, err := c.enc.Encode(v)
	if err != nil {
		return err
	}
	return c.carrier.Write(chunk)
}

func (c *Conn) Read(v interface{}) (err error) {
	chunk, err := c.carrier.Read()
	if err != nil && err != io.EOF {
		return err
	}
	return c.dec.Decode(chunk, v)
}

func (c *Conn) Close() (err error) {
	return c.carrier.Close()
}
