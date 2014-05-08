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

func cmdcmd(c proc.Cmd) Cmd {
	return Cmd{
		Env: c.Env,
		Path: c.Path,
		Args: c.Args,
	}
}

func (cmd Cmd) cmd() proc.Cmd {
	return proc.Cmd{
		Env: cmd.Env,
		Path: cmd.Path,
		Args: cmd.Args,
	}
}

type Stat struct {
	Cmd Cmd
	Exit error
	Phase string
}

func statstat(s proc.Stat) Stat {
	return Stat{
		Cmd: cmdcmd(s.Cmd),
		Exit: s.Exit,
		Phase: s.Phase,
	}
}

type Proc interface {
	Wait() (Stat, error)
	Signal(sig string) error
	GetEnv() []string
	GetCmd() Cmd
	IsDone() bool
	Peek() Stat
	Scrub()
}

type yprocProc struct {
	proc.YProc
}

func (y yprocProc) Wait() (Stat, error) {
	s, err := y.YValue.Wait()
	if err != nil {
		return Stat{}, err
	}
	return statstat(s), nil
}

func (y yprocProc) GetCmd() Cmd {
	return cmdcmd(y.YValue.GetCmd())
}

func (y yprocProc) Peek() Stat {
	return statstat(y.YValue.Peek())
}
