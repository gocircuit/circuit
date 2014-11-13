// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package dns

import (
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XNameserver{})
}

// X
type XNameserver struct {
	Nameserver
}

func (x XNameserver) Set(rr string) error {
	err := x.Nameserver.Set(rr)
	return errors.Pack(err)
}

// Y
type YNameserver struct {
	X circuit.X
}

func (y YNameserver) Set(rr string) error {
	r := y.X.Call("Set", rr)
	return errors.Unpack(r[0])
}

func (y YNameserver) Unset(name string) {
	y.X.Call("Unset", name)
}

func (y YNameserver) Scrub() {
	y.X.Call("Scrub")
}

func (y YNameserver) Peek() Stat {
	return y.X.Call("Peek")[0].(Stat)
}
