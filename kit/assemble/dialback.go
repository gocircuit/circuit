// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package assemble

import (
	"sync"

	"github.com/hoijui/circuit/use/circuit"
	"github.com/hoijui/circuit/use/n"
)

type DialBack struct {
	once sync.Once
	ch chan n.Addr
}

func NewDialBack() (*DialBack, *XDialBack) {
	d := &DialBack{ch: make(chan n.Addr, 1)}
	xd := &XDialBack{d}
	return d, xd
}

func (d *DialBack) ObtainAddr() n.Addr {
	return <-d.ch
}

type XDialBack struct {
	d *DialBack
}

func (xd *XDialBack) OfferAddr(addr n.Addr) {
	xd.d.once.Do(func() {
		xd.d.ch <- addr
		close(xd.d.ch)
	})
}

func init() {
	circuit.RegisterValue(&XDialBack{})
}

type YDialBack struct {
	circuit.PermX
}

func (y YDialBack) OfferAddr(addr n.Addr) {
	defer func() {
		recover()
	}()
	y.Call("OfferAddr", addr)
}
