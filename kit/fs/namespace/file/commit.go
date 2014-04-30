// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package file

import (
	"sync"

	"github.com/gocircuit/circuit/kit/fs/rh"
)

// CommitFile is an FID representing an open read-only file, 
// which acquires its underlying reader on the first call to Read or ReadIntr.
type CommitFile struct {
	*openFile
	slk  sync.Mutex
	src  rh.FID // CommitFile turns into a ReaderFile/WriterFile after the first read/write.
	clk  sync.Once
	cfn  CommitFunc
}

type CommitFunc func(rh.Intr) (rh.FID, error)

func NewCommitFile(cfn CommitFunc) rh.FID {
	return &CommitFile{
		openFile: newOpenFile(),
		cfn: cfn,
	}
}

func (f *CommitFile) source() rh.FID {
	f.slk.Lock()
	defer f.slk.Unlock()
	if f.src != nil {
		return f.src
	}
	return f.openFile
}

func (f *CommitFile) commit(intr rh.Intr) (fid rh.FID, err error) {
	f.clk.Do(
		func() {
			fid, err = f.cfn(intr)
			if err != nil {
				return
			}
			f.cfn = nil // good for gc
			f.slk.Lock()
			defer f.slk.Unlock()
			f.src = fid
		},
	)
	return f.source(), err
}

func (f *CommitFile) Read(at int64, count int, intr rh.Intr) (chunk rh.Chunk, err error) {
	fid, err := f.commit(intr)
	if err != nil {
		return nil, err
	}
	return fid.Read(at, count, intr)
}

func (f *CommitFile) Write(at int64, data rh.Chunk, intr rh.Intr) (n int, err error) {
	fid, err := f.commit(intr)
	if err != nil {
		return 0, err
	}
	return fid.Write(at, data, intr)
}

func (f *CommitFile) Clunk() (err error) {
	return f.source().Clunk()
}

func (f *CommitFile) Stat() (dir *rh.Dir, err error) {
	return f.source().Stat()
}
