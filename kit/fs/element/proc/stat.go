// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"bytes"

	"github.com/gocircuit/circuit/kit/iomisc"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

type StatFile struct {
	p *Proc
}

func NewStatFile(p *Proc) file.File {
	return &StatFile{p: p}
}

func (f *StatFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *StatFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.ReadOnly {
		return nil, rh.ErrPerm
	}
	return file.NewOpenReaderFile(iomisc.ReaderNopCloser(bytes.NewBufferString(f.p.GetState().String()))), nil
}

func (f *StatFile) Remove() error {
	return rh.ErrPerm
}
