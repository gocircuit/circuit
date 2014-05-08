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

func retypeProcStat(c proc.Cmd) Cmd {
	return Cmd{
		Env: c.Env,
		Path: c.Path,
		Args: c.Args,
	}
}

func (cmd Cmd) retype() proc.Cmd {
	return proc.Cmd{
		Env: cmd.Env,
		Path: cmd.Path,
		Args: cmd.Args,
	}
}

type ProcStat struct {
	Cmd Cmd
	Exit error
	Phase string
}

func statstat(s proc.Stat) ProcStat {
	return ProcStat{
		Cmd: retypeProcStat(s.Cmd),
		Exit: s.Exit,
		Phase: s.Phase,
	}
}

type Proc interface {
	Wait() (ProcStat, error)
	Signal(sig string) error
	GetEnv() []string
	GetCmd() Cmd
	IsDone() bool
	Peek() ProcStat
	Scrub()
}

type yprocProc struct {
	proc.YProc
}

func (y yprocProc) Wait() (ProcStat, error) {
	s, err := y.YProc.Wait()
	if err != nil {
		return ProcStat{}, err
	}
	return statstat(s), nil
}

func (y yprocProc) GetCmd() Cmd {
	return retypeProcStat(y.YProc.GetCmd())
}

func (y yprocProc) Peek() ProcStat {
	return statstat(y.YProc.Peek())
}
