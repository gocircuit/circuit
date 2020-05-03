// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package shell

import (
	"io"

	xyexec "github.com/hoijui/circuit/pkg/kit/x/exec"
	xyio "github.com/hoijui/circuit/pkg/kit/x/io"
	"github.com/hoijui/circuit/pkg/use/circuit"
	"github.com/hoijui/circuit/pkg/use/errors"
)

// XXX: Run commands in their own process groups?
// XXX: Add chrooting for all commands executed.

// XShell exports chroot-ed shells.
type XShell struct {
	chroot string
}

func init() {
	circuit.RegisterValue(&XShell{})
}

func NewXShell(chroot string) *XShell {
	return &XShell{chroot: chroot}
}

// Shell starts a new shell in a chroot and returns a cross-interface to its standard input, output and error.
func (ss *XShell) Shell() (xstdin circuit.X, xstdout, xstderr circuit.X, err error) {
	return xyexec.XCommand("sh", true).Start()
}

// ShellInteractive starts a new interactive shell in a chroot and returns a cross-interface to its pseudo-terminal.
func (ss *XShell) ShellInteractive() (xpty circuit.X, err error) {
	return xyexec.XCommand("sh", true, "-i").StartPTY([]string{"sane"})
}

/*
func (ss *XShell) ShellInteractive() (xstdin circuit.X, xstdout, xstderr circuit.X, err error) {
	return xyexec.XCommand("sh", true, "-i").Start()
}
*/

// Tail …
// XXX: Varargs doesn't work in cross-calls. Fix this.
func (ss *XShell) Tail(arg ...string) (xstdin circuit.X, xstdout, xstderr circuit.X, err error) {
	return xyexec.XCommand("tail", false, arg...).Start()
}

// YShell facilitates cross-calls to a XShell on the client side
type YShell struct {
	circuit.X
}

func (y YShell) Shell() (stdin io.WriteCloser, stdout, stderr io.ReadCloser, err error) {
	r := y.Call("Shell")
	err = errors.Unpack(r[3])
	if err != nil {
		return nil, nil, nil, err
	}
	return xyio.NewYWriteCloser(r[0]), xyio.NewYReadCloser(r[1]), xyio.NewYReadCloser(r[2]), nil
}

func (y YShell) Tail(name string, arg ...string) (stdin io.WriteCloser, stdout, stderr io.ReadCloser, err error) {
	ig := make([]interface{}, 0, len(arg))
	for _, a := range arg {
		ig = append(ig, a)
	}
	r := y.Call("Tail", ig...)
	err = errors.Unpack(r[3])
	if err != nil {
		return nil, nil, nil, err
	}
	return xyio.NewYWriteCloser(r[0]), xyio.NewYReadCloser(r[1]), xyio.NewYReadCloser(r[2]), nil
}
