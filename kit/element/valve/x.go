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
	Valve Valve
}

func (x XValve) Send() (circuit.X, error) {
	w, err := x.Valve.Send()
	if err != nil {
		return nil, err // errors created with errors.New are registered for cross-passing
	}
	return xio.NewXWriteCloser(w), nil
}

func (x XValve) Close() error {
	return x.Valve.Close()
}

func (x XValve) Recv() (circuit.X, error) {
	r, err := x.Valve.Recv()
	if err != nil {
		return nil, err
	}
	return xio.NewXReadCloser(r), nil
}

func (x XValve) Scrub() {
	x.Valve.Scrub()
}

func (x XValve) IsDone() bool {
	return x.Valve.IsDone()
}

func (x XValve) Cap() int {
	return x.Valve.Cap()
}

func (x XValve) Stat() *Stat {
	return x.Valve.Stat()
}

type YValve struct {
	X circuit.X
}

// all methods below will panic on system-level errors

func (y YValve) Send() (_ io.WriteCloser, err error) {
	r := y.X.Call("Send")
	if err = errors.Unpack(r[1]); err != nil {
		return nil, err
	}
	return xio.NewYWriteCloser(r[0]), nil
}

func (y YValve) Close() error {
	return errors.Unpack(y.X.Call("Close")[0])
}

func (y YValve) Recv() (_ io.ReadCloser, err error) {
	r := y.X.Call("Recv")
	if err = errors.Unpack(r[1]); err != nil {
		return nil, err
	}
	return xio.NewYReadCloser(r[0]), nil
}

func (y YValve) Cap() int {
	return y.X.Call("Cap")[0].(int)
}

func (y YValve) Stat() *Stat {
	return y.X.Call("Stat")[0].(*Stat)
}

func (y YValve) Scrub() {
	y.X.Call("Scrub")
}

func (y YValve) IsDone() bool {
	return y.X.Call("IsDone")[0].(bool)
}
