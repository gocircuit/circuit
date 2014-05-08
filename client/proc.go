// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"github.com/gocircuit/circuit/kit/element/proc"
)

type Cmd struct {
	Env []string
	Path string
	Args []string
}

type Stat struct {
	Cmd Cmd
	Exit error
	Phase string
}

type Proc interface {
	Wait() (Stat, error)
	Signal(sig string) error
	GetEnv() []string
	GetCmd() Cmd
	IsDone() bool
	Peek() Stat
}

type yprocProc struct {
	y proc.YProc
}

func (y yprocProc) Wait() (Stat, error) {
	??
}

func (y yprocProc) Signal(sig string) error {
	??
}

func (y yprocProc) GetEnv() []string {
	??
}

func (y yprocProc) GetCmd() Cmd {
	??
}

func (y yprocProc) IsDone() bool {
	??
}

func (y yprocProc) Peek() Stat {
	??
}

