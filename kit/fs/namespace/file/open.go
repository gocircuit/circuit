// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package file

import (
	"time"

	"github.com/gocircuit/circuit/kit/lang"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// openFile implements a subset of FID functionality, shared amongst open FIDs for reading/writing.
type openFile struct {
	rh.ZeroFID
	q rh.Q
	t time.Time
}

func newOpenFile() *openFile {
	f := &openFile{
		t: time.Now(),
	}
	f.q = rh.Q{
		ID:  uint64(lang.ComputeReceiverID(f)),
		Ver: 1,
	}
	return f
}

func (fid *openFile) String() string {
	return "open file"
}

func (fid *openFile) Q() rh.Q {
	return fid.q
}

func (fid *openFile) Walk(name []string) (rh.FID, error) {
	panic("open FID does not walk")
}

func (fid *openFile) Open(rh.Flag, rh.Intr) error {
	panic("FID already open")
}

func (fid *openFile) Clunk() error {
	return nil
}

func (fid *openFile) Stat() (*rh.Dir, error) {
	return &rh.Dir{
		Q:      fid.q,
		Mode:   rh.Mode{Attr: rh.ModeIO},
		Atime:  fid.t,
		Mtime:  fid.t,
		Name:   "",
		Length: 0,
		Uid:    rh.UID(),
		Gid:    rh.GID(),
	}, nil
}
