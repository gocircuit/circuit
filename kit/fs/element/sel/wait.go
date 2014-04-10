// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package sel

import (
	"bytes"

	"github.com/gocircuit/circuit/kit/iomisc"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

type WaitFile struct {
	s *Select
}

func NewWaitFile(s *Select) file.File {
	return &WaitFile{s: s}
}

func (f *WaitFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *WaitFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	f.s.ErrorFile.Set("") // clear error file
	if flag.Attr != rh.ReadOnly {
		return nil, rh.ErrPerm
	}
	waitfile, err := f.s.Wait(intr)
	if err == rh.ErrIntr {
		return nil, err
	}
	if waitfile != "" {
		return file.NewOpenReaderFile(iomisc.ReaderNopCloser(bytes.NewBufferString(waitfile))), nil
	}
	return nil, rh.ErrClash
}

func (f *WaitFile) Remove() error {
	return rh.ErrPerm
}
