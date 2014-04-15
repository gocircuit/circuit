// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"io"
)

type Chan interface {
	Send() io.WriteCloser
	Recv() io.ReadCloser
	Close()
	TrySend() io.WriteCloser
	TryRecv() io.ReadCloser
	SetCap(int)
	GetCap() int
}

type Command string {
	Env  []string
	Path string
	Args []string
}

type Proc interface {
	Start(Command) error
	Wait() error
}

type SelectClause interface {
	Send Chan
	Recv Chan
	Exit Proc
}

type Select interface {
	Start([]Clause)
	Wait() (int, error)
}
