// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package dir

import (
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// childFID wraps an FID and captures its Remove method, so that
// if the FID's Remove method returns no error on invocation,
// the FID's entry in the parent directory is removed.
type childFID struct {
	parent *Dir
	name   string
	rh.FID
}

func NewChildFID(parent *Dir, name string, fid rh.FID) *childFID {
	return &childFID{
		parent: parent,
		name:   name,
		FID:    fid,
	}
}

func (c *childFID) Walk(wname []string) (rh.FID, error) {
	if len(wname) == 0 { // cloning this fid clones the child wrapper and the wrapper fid
		walk, err := c.FID.Walk(nil)
		if err != nil {
			return nil, err
		}
		return &childFID{
			parent: c.parent,
			name:   c.name,
			FID:    walk,
		}, nil
	}
	return c.FID.Walk(wname)
}

func (c *childFID) Create(name string, flag rh.Flag, mode rh.Mode, perm rh.Perm) (rh.FID, error) {
	return c.FID.Create(name, flag, mode, perm)
}

func (c *childFID) Stat() (dir *rh.Dir, err error) {
	dir, err = c.FID.Stat()
	if err != nil {
		return nil, err
	}
	dir.Name = c.name
	return dir, nil
}

func (c *childFID) Remove() error {
	if err := c.FID.Remove(); err != nil {
		return err
	}
	if c.parent == nil {
		return rh.ErrPerm // cannot remove root
	}
	c.parent.RemoveChild(c.name)
	return nil
}
