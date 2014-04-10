// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package dash

import (
	"fmt"
	"path"

	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/element/valve"
	"github.com/gocircuit/circuit/kit/fs/element/proc"
	"github.com/gocircuit/circuit/kit/fs/element/sel"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

type NonRemovableDashDir struct {
	*DashDir
}

func (d NonRemovableDashDir) Remove() error {
	return rh.ErrPerm
}

// DashDir is a directory FID, constituting a dashboard (i.e. a workspace directory).
type DashDir struct {
	name string
	dir *dir.Dir
	rh.FID
	mkr []makerDir
}

type makerDir interface {
	rh.FID
	NumElements() int
}

func NewDir(name string) *DashDir {
	s := &DashDir{
		name: name,
		dir:  dir.NewDir(),
	}
	s.FID = s.dir.FID()
	s.dir.AddChild("help",
		file.NewFileFID(
			file.NewByteReaderFile(
				func() []byte {
					return []byte(s.help())
				},
			),
		),
	)
	s.mkr = []makerDir{
		valve.NewMakerDir(path.Join(s.name, "chan")),
		proc.NewMakerDir(path.Join(s.name, "proc")),
		sel.NewMakerDir(path.Join(s.name, "select")),
	}
	s.dir.AddChild("chan", s.mkr[0])
	s.dir.AddChild("proc", s.mkr[1])
	s.dir.AddChild("select", s.mkr[2])
	s.dir.AllowCreate()
	return s
}

func (s *DashDir) Walk(wname []string) (rh.FID, error) {
	if len(wname) > 0 {
		return s.FID.Walk(wname)
	}
	return s, nil
}

func (s *DashDir) Create(name string, flag rh.Flag, mode rh.Mode, perm rh.Perm) (rh.FID, error) {
	if mode.Attr != rh.ModeDir {
		return nil, rh.ErrPerm
	}
	return s.dir.AddChild(name, NewDir(name))
}

func (s *DashDir) Remove() error {
	s.dir.DisallowCreate()
	// check for descendant elements
	for _, d := range s.mkr {
		if d.NumElements() > 0 {
			s.dir.AllowCreate()
			return rh.ErrBusy
		}
	}
	// check for child dashboards
	n := s.dir.Len() - len(s.mkr) - 1 // sub-dashes = children - maker directories and help file
	if n > 0 {
		s.dir.AllowCreate()
		return rh.ErrBusy
	}
	return nil
}

func (s *DashDir) help() string {
	return fmt.Sprintf(helpFormat, s.name)
}

const helpFormat = `
	This is a circuit dashboard directory named: %s

	A dashboard is a working space wherein the user can create
	and manipulate circuit elements (channels, processes,
	selectors, etc).

	To learn about and to create specific element types, investigate the 
	subdirectories of this dashboard. For instance, to work with channels:

		cd chan
		cat help

MKDIR

	Directories created within this directory automatically become
	dashboard directories. For instance,

		mkdir subdash

	Dashboard directories are effectively namespaces within which
	the user can create multiple circuit elements.

RMDIR

	Subordinate dashboard directories can be removed with "rmdir"
	as long as (i) they don't have subordinate dashboards themselves
	and (ii) their element directories (chan, proc and select) do
	not have any user elements inside them.

		rmdir subdash

	The "help" file as well as the element subdirectories "chan",
	"proc" and "select" of this dashboard cannot be removed.
`
