// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"io"
	"runtime"

	"github.com/gocircuit/circuit/element/proc"
	"github.com/gocircuit/circuit/kit/interruptible"
	"github.com/gocircuit/circuit/kit/lang"
	"github.com/gocircuit/circuit/use/circuit"

	dkr "github.com/fsouza/go-dockerclient"
)

type Container interface {
	Scrub()
	Wait() (*dkr.Container, error)
	Signal(sig string) error
	IsDone() bool
	Peek() (*dkr.Container, error)
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
	X() circuit.X
}

type container struct {
	dkr *dkr.Client
	name string
	con *dkr.Container
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	exit <-chan error
}

func MakeContainer(run Run) (_ Container, err error) {
	con := &container{
		name: "via-circuit-"+lang.ChooseReceiverID().String()[1:],
	}
	if con.dkr, err = dial(); err != nil {
		return nil, err
	}
	var cmd []string
	if run.Path != "" {
		cmd = append([]string{run.Path}, run.Args...)
	}
	opts := dkr.CreateContainerOptions{
		Name: con.name,
		Config: &dkr.Config{
			Memory: run.Memory,
			CpuShares: run.CpuShares,
			AttachStdin: true,
			AttachStdout: true,
			AttachStderr: true,
			// OpenStdin: xxx,
			// StdinOnce: xxx,
			Env: run.Env,
			Cmd: cmd,
			Image: run.Image,
			Volumes: run.Volumes(),
			WorkingDir: run.Dir,
    		},
	}
	con.con, err = con.dkr.CreateContainer(opts)
	if err != nil {
		return nil, err
	}
	// Attach standard streams and log stream
	if err = con.attach(); err != nil {
		return nil, err
	}
	// Exit mechanism
	ch := make(chan error, 1)
	con.exit = ch
	go con.wait(ch)
	// Removal mechanism
	runtime.SetFinalizer(con,
		func(c *container) {
			opts := dkr.RemoveContainerOptions{
				ID: c.con.ID,
				// RemoveVolumes: xx,
				// Force: xx,
			}
			c.dkr.RemoveContainer(opts)
		},
	)
	return con, nil
}

func (con *container) Wait() (_ *dkr.Container, err error) {
	con.waitexit()
	return con.Peek()
}

func (con *container) waitexit() {
	<-con.exit
}

func (con *container) wait(ch chan<- error) {
	_, err := con.dkr.WaitContainer(con.con.ID)
	ch <- err
	close(ch)
	return
}

func (con *container) attach() error {
	var opts = dkr.AttachToContainerOptions{
		    Container: con.con.ID,
		    Stream: true,
		    Stdin: true,
		    Stdout: true,
		    Stderr: true,
	}
	const bufferLen = 32e3
	opts.InputStream, con.stdin = interruptible.BufferPipe(bufferLen)
	con.stdout, opts.OutputStream = interruptible.BufferPipe(bufferLen)
	con.stderr, opts.ErrorStream = interruptible.BufferPipe(bufferLen)
	return con.dkr.AttachToContainer(opts)
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

func (con *container) Peek() (*dkr.Container, error) {
	return con.dkr.InspectContainer(con.con.ID)
}

func (con *container) Scrub() {
	opts := dkr.RemoveContainerOptions{
		ID: con.con.ID,
		// RemoveVolumes: xx,
		// Force: xx,
	}
	con.dkr.RemoveContainer(opts)
}

func (con *container) Signal(sig string) error {
	signo, _ := proc.ParseSignal(sig)
	opts := dkr.KillContainerOptions{
		ID: con.con.ID,
		Signal: dkr.Signal(signo),
	}
	return con.dkr.KillContainer(opts)
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
	return nil
}
