// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package file

import (
	"sync"
	"time"

	"github.com/gocircuit/circuit/kit/lang"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

type File interface {

	// Perm returns the current permissions of this file.
	Perm() rh.Perm

	// Open returns an FID representing an open version of this file.
	Open(rh.Flag, rh.Intr) (rh.FID, error)

	// Remove is invoked when the file system user attempts to remove this file.
	// The operation will succeed at POSIX level, only if Remove returns nil.
	Remove() error
}

// file is an FID for an unopened file backed by a File interface.
type file struct {
	rh.ZeroFID
	q rh.Q
	t time.Time
	f File
	o struct{
		sync.Mutex
		rh.FID // FID of the opened file, if the user opens this file
	}
}

// NewFileFID returns a new FID for an unopened file, given by the File interface.
func NewFileFID(f File) rh.FID {
	return &file{
		q: rh.Q{
			ID:  uint64(lang.ComputeReceiverID(f)),
			Ver: 1,
		},
		t: time.Now(),
		f: f,
	}
}

func (fid *file) String() string {
	return "file"
}

func (fid *file) Q() rh.Q {
	return fid.q
}

func (fid *file) Walk(name []string) (rh.FID, error) {
	if len(name) > 0 {
		return nil, rh.ErrClash
	}
	return &file{
		q: fid.q,
		t: fid.t,
		f: fid.f,
	}, nil
}

func (fid *file) Open(flag rh.Flag, intr rh.Intr) (err error) {
	fid.o.Lock()
	defer fid.o.Unlock()
	if fid.o.FID != nil {
		return rh.ErrBusy
	}
	fid.o.FID, err = fid.f.Open(flag, intr)
	return
}

func (fid *file) Wstat(wdir *rh.Wdir) error {
	return nil
}

func (fid *file) Read(offset int64, count int, intr rh.Intr) (rh.Chunk, error) {
	fid.o.Lock()
	defer fid.o.Unlock()
	if fid.o.FID == nil {
		return nil, rh.ErrClash
	}
	return fid.o.Read(offset, count, intr)
}

func (fid *file) Write(offset int64, data rh.Chunk, intr rh.Intr) (int, error) {
	fid.o.Lock()
	defer fid.o.Unlock()
	if fid.o.FID == nil {
		return 0, rh.ErrClash
	}
	return fid.o.Write(offset, data, intr)
}

func (fid *file) Clunk() (err error) {
	fid.o.Lock()
	defer fid.o.Unlock()
	if fid.o.FID == nil {
		return rh.ErrGone
	}
	err = fid.o.FID.Clunk()
	fid.o.FID = nil
	return
}

func (fid *file) Stat() (*rh.Dir, error) {
	return &rh.Dir{
		Q:      fid.q,
		Mode:   rh.Mode{Attr: rh.ModeIO},
		Atime:  fid.t,
		Mtime:  fid.t,
		Name:   "",
		Perm:   fid.f.Perm(),
		Length: 0,
		Uid:    rh.UID(),
		Gid:    rh.GID(),
	}, nil
}
