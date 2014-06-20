// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"io"
)

type Container interface {
	Scrub()
	IsDone() bool
	Peek() (*Stat, error)
	Signal(sig string) error
	Wait() (*Stat, error)
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
}
