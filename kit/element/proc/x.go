// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XProc{})
}

type XProc struct {
	Proc
}

func (x XProc) Wait() (Stat, error) {
	stat, err := x.Proc.Wait()
	return stat, errors.Pack(err)
}

func (x XProc) Signal(sig string) error {
	return errors.Pack(x.Proc.Signal(sig))
}

type YProc struct {
	X circuit.X
}

func (y YProc) Wait() (Stat, error) {
	r := y.X.Call("Wait")
	return r[0].(Stat), errors.Unpack(r[1])
}

func (y YProc) Signal(sig string) error {
	r := y.X.Call("Signal")
	return errors.Unpack(r[0])
}

func (y YProc) Scrub() {
	y.X.Call("Scrub")
}

func (y YProc) GetEnv() []string {
	return y.X.Call("GetEnv")[0].([]string)
}

func (y YProc) GetCmd() Cmd {
	return y.X.Call("GetCmd")[0].(Cmd)
}

func (y YProc) IsDone() bool {
	return y.X.Call("IsDone")[0].(bool)
}

func (y YProc) Peek() Stat {
	return y.X.Call("Peek")[0].(Stat)
}
