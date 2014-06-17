// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"io"

	"github.com/gocircuit/circuit/use/circuit"
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
}

func MakeContainer(run Run) Container {
	panic(1)
}
