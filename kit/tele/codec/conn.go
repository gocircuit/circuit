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

type Conn struct {
	enc   Encoder
	dec   Decoder
	faith *faithful.Conn
}

func NewConn(faith *faithful.Conn, codec Codec) *Conn {
	return &Conn{
		enc:   codec.NewEncoder(),
		dec:   codec.NewDecoder(),
		faith: faith,
	}
}

func (c *Conn) String() string {
	return c.RemoteAddr().String()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.faith.RemoteAddr()
}

func (c *Conn) Write(v interface{}) error {
	chunk, err := c.enc.Encode(v)
	if err != nil {
		return err
	}
	return c.faith.Write(chunk)
}

func (c *Conn) Read(v interface{}) error {
	chunk, err := c.faith.Read()
	if err != nil {
		return err
	}
	return c.dec.Decode(chunk, v)
}

func (c *Conn) Close() error {
	return c.faith.Close()
}
