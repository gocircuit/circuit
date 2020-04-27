// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package exec facilitates remote execution of OS processes across workers
package exec

import (
	"io"
	"os/exec"
	"sync"
	"syscall"

	"github.com/hoijui/circuit/kit/pty"
	xyio "github.com/hoijui/circuit/kit/x/io"
	"github.com/hoijui/circuit/use/circuit"
	"github.com/hoijui/circuit/use/errors"
)

// XCmd exports chroot-ed shells.
type XCmd struct {
	cmd *exec.Cmd
}

func init() {
	circuit.RegisterValue(&XCmd{})
}

func XCommand(name string, setsid bool, arg ...string) *XCmd {
	x := &XCmd{cmd: exec.Command(name, arg...)}
	if setsid {
		x.cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	}
	return x
}

func (x *XCmd) StartPTY(stty []string) (xpty circuit.X, err error) {
	pty, tty, err := pty.Open()
	if err != nil {
		return nil, errors.Pack(err)
	}
	defer tty.Close() // Close our handle on the TTY after the child executes
	x.cmd.Stdout, x.cmd.Stdin, x.cmd.Stderr = tty, tty, tty
	x.cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	if err = x.cmd.Start(); err != nil {
		pty.Close()
		return nil, errors.Pack(err)
	}
	// sfs1/rex has comment that maybe tty flags cannot change until tty opened by child?? Not sure.
	if len(stty) > 0 {
		arg0 := []string{"-f", tty.Name()}
		if err = exec.Command("stty", append(arg0, "sane")...).Run(); err != nil {
			return nil, errors.Pack(err)
		}
		//log.Printf("PTY=%s, XX %#v", pty.Name(), append(arg0, stty...))
		if err = exec.Command("stty", append(arg0, stty...)...).Run(); err != nil {
			return nil, errors.Pack(err)
		}
	}
	defer func() {
		go x.cmd.Wait() // Calling wait is necessary, otherwise child remains zombie
	}()
	return xyio.NewXReadWriteCloser(pty), nil
}

// Start starts the command and returns cross-interfaces to its standard input, output and error.
func (x *XCmd) Start() (xstdin circuit.X, xstdout, xstderr circuit.X, err error) {
	//
	stdin, err := x.cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	//
	stdout, err := x.cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	//
	stderr, err := x.cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	//
	if err = x.cmd.Start(); err != nil {
		return nil, nil, nil, err
	}
	defer func() {
		go x.cmd.Wait() // Calling wait is necessary, otherwise child remains zombie
	}()
	return xyio.NewXWriteCloser(stdin), xyio.NewXReadCloser(stdout), xyio.NewXReadCloser(stderr), nil
}

// Wait blocks until the process completes, akin to os/exec.Cmd.Wait
func (x *XCmd) Wait() error {
	return errors.Pack(x.cmd.Wait())
}

// Run executes the command and blocks until it completes, akin to os/exec.Cmd.Run
func (x *XCmd) Run() error {
	return errors.Pack(x.cmd.Run())
}

// onRC ensures that a custom onclose function is invoked after the underlying I/O object is closed.
type onRC struct {
	io.ReadCloser
	sync.Once
	onclose func()
}

func (z *onRC) Close() error {
	defer z.Once.Do(z.onclose)
	return z.ReadCloser.Close()
}

// onWC ensures that a custom onclose function is invoked after the underlying I/O object is closed.
type onWC struct {
	io.WriteCloser
	sync.Once
	onclose func()
}

func (z *onWC) Close() error {
	defer z.Once.Do(z.onclose)
	return z.WriteCloser.Close()
}

// YCmd facilitates usage of a cross-interface to XCmd on the caller side
type YCmd struct {
	circuit.X // Cross-interface to *XCmd
}

func (y YCmd) Start() (stdin io.WriteCloser, stdout, stderr io.ReadCloser, err error) {
	r := y.Call("Start")
	err = errors.Unpack(r[3])
	if err != nil {
		return nil, nil, nil, err
	}
	return xyio.NewYWriteCloser(r[0]), xyio.NewYReadCloser(r[1]), xyio.NewYReadCloser(r[2]), nil
}

func (y YCmd) Wait() error {
	return errors.Unpack(y.Call("Wait"))
}

func (y YCmd) Run() error {
	return errors.Unpack(y.Call("Run"))
}
