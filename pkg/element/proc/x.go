// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"io"
	
	xio "github.com/hoijui/circuit/pkg/kit/x/io"
	"github.com/hoijui/circuit/pkg/use/circuit"
	"github.com/hoijui/circuit/pkg/use/errors"
)

func init() {
	circuit.RegisterValue(XProc{})
}

type XProc struct {
	Proc
}

func unpack(stat Stat) Stat {
	stat.Exit = errors.Unpack(stat.Exit)
	return stat
}

func pack(stat Stat) Stat {
	stat.Exit = errors.Pack(stat.Exit)
	return stat
}

func (x XProc) Wait() (Stat, error) {
	stat, err := x.Proc.Wait()
	return pack(stat), errors.Pack(err)
}

func (x XProc) Signal(sig string) error {
	return errors.Pack(x.Proc.Signal(sig))
}

func (x XProc) Stdin() circuit.X {
	return xio.NewXWriteCloser(x.Proc.Stdin())
}

func (x XProc) Stdout() circuit.X {
	return xio.NewXReadCloser(x.Proc.Stdout())
}

func (x XProc) Stderr() circuit.X {
	return xio.NewXReadCloser(x.Proc.Stderr())
}

func (x XProc) Peek() Stat {
	return pack(x.Proc.Peek())
}

type YProc struct {
	X circuit.X
}

func (y YProc) Wait() (Stat, error) {
	r := y.X.Call("Wait")
	return unpack(r[0].(Stat)), errors.Unpack(r[1])
}

func (y YProc) Signal(sig string) error {
	r := y.X.Call("Signal", sig)
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
	return unpack(y.X.Call("Peek")[0].(Stat))
}

func (y YProc) Stdin() io.WriteCloser {
	return xio.NewYWriteCloser(y.X.Call("Stdin")[0])
}

func (y YProc) Stdout() io.ReadCloser {
	return xio.NewYReadCloser(y.X.Call("Stdout")[0])
}

func (y YProc) Stderr() io.ReadCloser {
	return xio.NewYReadCloser(y.X.Call("Stderr")[0])
}
