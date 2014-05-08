// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"io"

	"github.com/gocircuit/circuit/kit/element/valve"
)

type Chan interface {
	Send() (io.WriteCloser, error)
	IsDone() bool
	Scrub()
	Close() error
	Recv() (io.ReadCloser, error)
	Cap() int
	Stat() Stat
}

type Stat struct {
	Cap int
	Opened bool
	Closed bool
	Aborted bool
	NumSend int
	NumRecv int
}

type yvalveChan struct {
	y valve.YValue
}

func (y yvalveChan) Send() (io.WritCloser, error) {
	??
}

func (y yvalveChan) IsDone() bool {
	??
}

func (y yvalveChan) Scrub() {
	??
}

func (y yvalveChan) Close() error {
	??
}

func (y yvalveChan) Recv() (io.ReadCloser, error) {
	??
}

func (y yvalveChan) Cap() int {
	??
}

func (y yvalveChan) Stat() Stat {
	??
}

