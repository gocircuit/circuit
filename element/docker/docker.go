// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"errors"
	"io"
	"sync"

	"github.com/gocircuit/circuit/kit/lang"
	"github.com/gocircuit/circuit/use/circuit"

	dkr "github.com/fsouza/go-dockerclient"
)

type Container interface {
	Scrub()
	Wait() (Stat, error)
	Signal(sig string) error
	IsDone() bool
	Peek() Stat
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
	X() circuit.X
}

type container struct {
	id lang.ReceiverID
	sys struct {
		sync.Mutex
		con *dkr.Container
		exit error
	}
}

func MakeContainer(run Run) (_ *container, err error) {
	cli := client()
	if cli == nil {
		return nil, errors.New("docker not enabled")
	}
	con := &container{
		id: lang.ChooseReceiverID(),
	}
	var cmd []string
	if run.Path != "" {
		cmd = append([]string{run.Path}, run.Args...)
	}
	opts := dkr.CreateContainerOptions{
		Name: con.id.String(),
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
	con.sys.con, err = cli.CreateContainer(opts)
	if err != nil {
		return nil, err
	}
	return con, nil
}
