// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"

	"github.com/gocircuit/circuit/element/proc"
	"github.com/gocircuit/circuit/kit/interruptible"
	"github.com/gocircuit/circuit/kit/lang"
	"github.com/gocircuit/circuit/use/circuit"
	ds "github.com/gocircuit/circuit/client/docker"
)

type Container interface {
	Scrub()
	IsDone() bool
	Peek() (*ds.Stat, error)
	Signal(sig string) error
	Wait() (*ds.Stat, error)
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
	X() circuit.X
}

type container struct {
	name string
	cmd *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	exit <-chan error
}

func MakeContainer(run ds.Run) (_ Container, err error) {
	ch := make(chan error, 1)
	con := &container{
		name: "via-circuit-"+lang.ChooseReceiverID().String()[1:],
		exit: ch,
	}
	con.cmd = exec.Command(dkr, run.Arg(con.name)...)
	println(fmt.Sprintf("%v", run.Arg(con.name)))
	con.cmd.Stdin, con.stdin = interruptible.BufferPipe(StdBufferLen)
	con.stdout, con.cmd.Stdout = interruptible.BufferPipe(StdBufferLen)
	con.stderr, con.cmd.Stderr = interruptible.BufferPipe(StdBufferLen)
	if err = con.cmd.Start(); err != nil {
		return nil, err
	}
	go func() {
		ch <- con.cmd.Wait()
		close(ch)
		con.stdout.Close()
		con.stderr.Close()
	}()
	runtime.SetFinalizer(con,
		func(c *container) {
			exec.Command(dkr, "rm", c.name).Run()
		},
	)
	return con, nil
}

func (con *container) Wait() (_ *ds.Stat, err error) {
	<-con.exit
	return con.Peek()
}

func (con *container) Stdin() io.WriteCloser {
	return con.stdin
}

func (con *container) Stdout() io.ReadCloser {
	return con.stdout
}

func (con *container) Stderr() io.ReadCloser {
	return con.stderr
}

func (con *container) Peek() (stat *ds.Stat, err error) {
	buf, err := exec.Command(dkr, "inspect", con.name).Output()
	if err != nil {
		return nil, err
	}
	if stat, err = ds.ParseStatInArray(buf); err != nil {
		return nil, err
	}
	return
}

func (con *container) Scrub() {
	exec.Command(dkr, "rm", con.name).Run()
}

func (con *container) Signal(sig string) error {
	signo, ok := proc.ParseSignal(sig)
	if !ok {
		return errors.New("signal name not recognized")
	}
	if con.cmd.Process == nil {
		return errors.New("no process")
	}
	return con.cmd.Process.Signal(signo)
}

func (con *container) IsDone() bool {
	select {
	case <-con.exit:
		return true
	default:
		return false
	}
}

func (con *container) X() circuit.X {
	panic(0)
}
