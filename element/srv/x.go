// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package srv

import (
	"io"

	xio "github.com/gocircuit/circuit/kit/x/io"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XServer{})
}

// XServer…
type XServer struct {
	*server
}

func (x XServer) Profile(name string) (circuit.X, error) {
	r, err := x.server.Profile(name)
	if err != nil {
		return nil, errors.Pack(err)
	}
	return xio.NewXReader(r), nil
}

// YServer…
type YServer struct {
	X circuit.X
}

func (y YServer) Profile(name string) (io.Reader, error) {
	r := y.X.Call("Profile", name)
	if err := errors.Unpack(r[1]); err != nil {
		return nil, err
	}
	return xio.NewYReader(r[0]), nil
}

func (y YServer) Peek() Stat {
	return y.X.Call("Peek")[0].(Stat)
}

func (y YServer) IsDone() bool {
	return y.X.Call("IsDone")[0].(bool)
}

func (y YServer) Scrub() {
	y.X.Call("Scrub")
}
