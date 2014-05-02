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
	circuit.RegisterValue(XValve{})
}

type XValve struct {
	v *Valve
}

func (x XValve) Send() (circuit.X, error) {
	w, err := x.v.Send()
	if err != nil {
		return nil, err // errors created with errors.New are registered for cross-passing
	}
	return xio.NewXWriteCloser(w), nil
}

func (x XValve) Close() error {
	return x.v.Close()
}

func (x XValve) Recv() (circuit.X, error) {
	r, err := x.v.Recv()
	if err != nil {
		return nil, err
	}
	return xio.NewXReadCloser(r), nil
}

func (x XValve) Cap() int {
	return x.v.Cap()
}

func (x XValve) Stat() *Stat {
	return x.v.Stat()
}

type YValve struct {
	x circuit.X
}

// all methods below will panic on system-level errors

func (y YValve) Send() (_ io.WriteCloser, err error) {
	r := y.x.Call("Send")
	if err = errors.Unpack(r[1]); err != nil {
		return nil, err
	}
	return xio.NewYWriteCloser(r[0]), nil
}

func (y YValve) Close() error {
	return errors.Unpack(y.x.Call("Close")[0])
}

func (y YValve) Recv() (_ io.ReadCloser, err error) {
	r := y.x.Call("Recv")
	if err = errors.Unpack(r[1]); err != nil {
		return nil, err
	}
	return xio.NewYReadCloser(r[0]), nil
}

func (y YValve) Cap() int {
	return y.x.Call("Cap")[0].(int)
}

func (y YValve) Stat() *Stat {
	return y.x.Call("Stat")[0].(*Stat)
}
