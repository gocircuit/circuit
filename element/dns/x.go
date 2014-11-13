// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package dns

import (
	"io"
	
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XNameserver{})
}

type XNameserver struct {
	Nameserver
}

func unpack(stat Stat) Stat {
	stat.Exit = errors.Unpack(stat.Exit)
	return stat
}

func pack(stat Stat) Stat {
	stat.Exit = errors.Pack(stat.Exit)
	return stat
}

func (x XNameserver) Wait() (Stat, error) {
	stat, err := x.Nameserver.Wait()
	return pack(stat), errors.Pack(err)
}

func (x XNameserver) Signal(sig string) error {
	return errors.Pack(x.Nameserver.Signal(sig))
}

func (x XNameserver) Stdin() circuit.X {
	return xio.NewXWriteCloser(x.Nameserver.Stdin())
}

func (x XNameserver) Stdout() circuit.X {
	return xio.NewXReadCloser(x.Nameserver.Stdout())
}

func (x XNameserver) Stderr() circuit.X {
	return xio.NewXReadCloser(x.Nameserver.Stderr())
}

func (x XNameserver) Peek() Stat {
	return pack(x.Nameserver.Peek())
}

type YNameserver struct {
	X circuit.X
}

func (y YNameserver) Wait() (Stat, error) {
	r := y.X.Call("Wait")
	return unpack(r[0].(Stat)), errors.Unpack(r[1])
}

func (y YNameserver) Signal(sig string) error {
	r := y.X.Call("Signal", sig)
	return errors.Unpack(r[0])
}

func (y YNameserver) Scrub() {
	y.X.Call("Scrub")
}

func (y YNameserver) GetEnv() []string {
	return y.X.Call("GetEnv")[0].([]string)
}

func (y YNameserver) GetCmd() Cmd {
	return y.X.Call("GetCmd")[0].(Cmd)
}

func (y YNameserver) IsDone() bool {
	return y.X.Call("IsDone")[0].(bool)
}

func (y YNameserver) Peek() Stat {
	return unpack(y.X.Call("Peek")[0].(Stat))
}

func (y YNameserver) Stdin() io.WriteCloser {
	return xio.NewYWriteCloser(y.X.Call("Stdin")[0])
}

func (y YNameserver) Stdout() io.ReadCloser {
	return xio.NewYReadCloser(y.X.Call("Stdout")[0])
}

func (y YNameserver) Stderr() io.ReadCloser {
	return xio.NewYReadCloser(y.X.Call("Stderr")[0])
}
