// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"io"

	xio "github.com/gocircuit/circuit/kit/x/io"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XProc{})
}

type XProc struct {
	p *Proc
}

func (x XProc) Send() (circuit.X, error) {
	w, err := x.v.Send()
	if err != nil {
		return nil, err // errors created with errors.New are registered for cross-passing
	}
	return xio.NewXWriteCloser(w), nil
}

func (p *Proc) Scrub() {
	???
}

type YProc struct {
	x circuit.X
}

// all methods below will panic on system-level errors

func (y YProc) Send() (_ io.WriteCloser, err error) {
	r := y.x.Call("Send")
	if err = errors.Unpack(r[1]); err != nil {
		return nil, err
	}
	return xio.NewYWriteCloser(r[0]), nil
}
