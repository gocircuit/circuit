// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package sel

import (
	"path"

	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// SelectMakerDir is an FID for a directory wherein the user can create valves
type SelectMakerDir struct {
	name string
	rh.FID
	dir *dir.Dir
}

func NewMakerDir(name string) *SelectMakerDir {
	d := &SelectMakerDir{
		name: name,
		dir:  dir.NewDir(),
	}
	d.FID = d.dir.FID()
	d.dir.AddChild("help", file.NewFileFID(file.NewByteReaderFile(
		func() []byte { return []byte(d.Help()) },
	)))
	d.dir.AllowCreate()
	return d
}

func (s *SelectMakerDir) Walk(wname []string) (rh.FID, error) {
	if len(wname) > 0 {
		return s.FID.Walk(wname)
	}
	return s, nil
}

func (s *SelectMakerDir) Create(name string, _ rh.Flag, mode rh.Mode, _ rh.Perm) (rh.FID, error) {
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

func (s *SelectMakerDir) Remove() error {
	return rh.ErrPerm // cannot remove individual maker directories
}

func (s *SelectMakerDir) NumElements() int {
	return s.dir.Len() - 1 // account for help file
}

func (s *SelectMakerDir) Help() string {
	return mkdirHelp
}

const mkdirHelp = `
	In this directory you can create circuit select elements.

	Select elements are analogous to select statements in the Go Language.
	They are a way of waiting until one of a collection of files
	is ready for reading.

	To make a new select element, make a subdirectory.

		mkdir sarah

	Navigate there and look around.

		cd sarah
		cat help

RMDIR

	Circuit select element directories can be removed with "rmdir",
	as long as a selection is not in progress.

		rmdir sarah

`
