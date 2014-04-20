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

// Channel

type Chan interface {
	Send() io.WriteCloser
	Close()
	Recv() io.ReadCloser
}

// Process

type Command struct {
	Env  []string `json:"env"`
	Path string   `json:"path"`
	Args []string `json:"args"`
}

type Proc interface {
	Start(Command) error
	Wait() error
}

// Select

type Clause interface{}

type ClauseDefault struct{}

type ClauseSend struct {
	Chan
}

type ClauseRecv struct {
	Chan
}

type ClauseExit struct {
	Proc
}

type Select interface {

	// Wait blocks until one of the select branches unblocks.
	// branch indicates the index of the clause that unblocked.
	// 
	Wait() (branch int, value interface{}) // io.WriteCloser, io.ReadCloser, error
}
