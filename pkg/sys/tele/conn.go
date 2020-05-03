// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tele

import (
	"github.com/hoijui/circuit/pkg/kit/tele/blend"
	"github.com/hoijui/circuit/pkg/use/n"
)

type Conn struct {
	addr *Addr
	sub  *blend.Conn
}

func NewConn(sub *blend.Conn, addr *Addr) *Conn {
	return &Conn{addr: addr, sub: sub}
}

func (c *Conn) Read() (v interface{}, err error) {
	if v, err = c.sub.Read(); err != nil {
		return nil, err
	}
	return
}

func (c *Conn) Write(v interface{}) (err error) {
	if err = c.sub.Write(v); err != nil {
		return err
	}
	return nil
}

func (c *Conn) Close() error {
	return c.sub.Close()
}

func (c *Conn) Abort(reason error) {
	c.sub.Abort(reason)
}

func (c *Conn) Addr() n.Addr {
	return c.addr
}
