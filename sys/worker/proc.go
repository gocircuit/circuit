// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package worker

import (
	"io"

	"github.com/gocircuit/circuit/use/n"
)

type Console struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

type Process struct {
	console Console
	addr    n.Addr
}

func (p *Process) Addr() n.Addr {
	return p.addr
}

func (p *Process) Kill() error {
	return kill(p.addr)
}

func (p *Process) Stdin() io.WriteCloser {
	panic("ni")
	return p.console.stdin
}

func (p *Process) Stdout() io.ReadCloser {
	panic("ni")
	return p.console.stdout
}

func (p *Process) Stderr() io.ReadCloser {
	panic("ni")
	return p.console.stderr
}
