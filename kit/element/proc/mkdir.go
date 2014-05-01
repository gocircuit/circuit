// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"path"

	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// ProcMakerDir is an FID for a directory wherein the user can create valves
type ProcMakerDir struct {
	name string
	rh.FID
	dir *dir.Dir
}

func NewMakerDir(name string) *ProcMakerDir {
	d := &ProcMakerDir{
		name: name,
		dir:  dir.NewDir(),
	}
	d.FID = d.dir.FID()
	d.dir.AddChild("help", file.NewFileFID(file.NewByteReaderFile(
		func() []byte {
			return []byte(d.Help())
		}),
	))
	d.dir.AllowCreate()
	return d
}

func (s *ProcMakerDir) Walk(wname []string) (rh.FID, error) {
	if len(wname) > 0 {
		return s.FID.Walk(wname)
	}
	return s, nil
}

func (s *ProcMakerDir) Create(name string, _ rh.Flag, mode rh.Mode, _ rh.Perm) (rh.FID, error) {
	if mode.Attr != rh.ModeDir {
		return nil, rh.ErrPerm
	}
	return s.dir.AddChild(
		name, 
		NewDir(
			path.Join(s.name, name), 
			func() {
				s.dir.RemoveChild(name)
			},
		),
	)
}

func (s *ProcMakerDir) Remove() error {
	return rh.ErrPerm // maker directories are themselves not removable from their dash
}

func (s *ProcMakerDir) NumElements() int {
	return s.dir.Len() - 1 // account for help file
}

func (s *ProcMakerDir) Help() string {
	return mkdirHelp
}

const mkdirHelp = `
	In this directory you can create and operate circuit processes.

	Circuit processes are a way of executing, managing and 
	synchronizing with conventional OS processes.

	To prepare a new process for execution, make a new directory.

		mkdir paul

	To configure, run and learn more about the newly created
	circuit process, navigate to its directory and look around.

		cd paul
		cat help

RMDIR

	Circuit process directories can be removed with "rmdir" as long
	as the underlying process has exited.

		rmdir paul

`
