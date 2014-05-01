// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"path"

	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// ValveMakerDir is an FID for a directory wherein the user can create valves
type ValveMakerDir struct {
	name string
	rh.FID
	dir *dir.Dir
}

func NewMakerDir(name string) *ValveMakerDir {
	d := &ValveMakerDir{
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

func (s *ValveMakerDir) Walk(wname []string) (rh.FID, error) {
	if len(wname) > 0 {
		return s.FID.Walk(wname)
	}
	return s, nil
}

func (s *ValveMakerDir) Create(name string, _ rh.Flag, mode rh.Mode, _ rh.Perm) (rh.FID, error) {
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

func (s *ValveMakerDir) Remove() error {
	return rh.ErrPerm // cannot remove individual maker directories
}

func (s *ValveMakerDir) NumElements() int {
	return s.dir.Len() - 1 // account for help file
}

func (s *ValveMakerDir) Help() string {
	return mkdirHelp
}

const mkdirHelp = `
	In this directory you can create and use circuit channels.

	Circuit channels are analogous to Go Language channels:
	They are a data structure for synchronization between a
	sender and a receiver, via send, receive and close operations.

	To create a new channel, make a subdirectory.

		mkdir charlie

	To use and learn more about the newly created channel,
	navigate to its directory and look around.

		cd charlie
		cat help

REMOVAL

	Circuit channel directories can be removed with "rmdir" only 
	after the channel has been closed.

		rmdir charlie

`
