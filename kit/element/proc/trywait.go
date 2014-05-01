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

type TryWaitFile struct {
	p *Proc
}

func NewTryWaitFile(p *Proc) file.File {
	return &TryWaitFile{p: p}
}

func (f *TryWaitFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *TryWaitFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.ReadOnly {
		return nil, rh.ErrPerm
	}
	return file.NewOpenReaderFile(iomisc.ReaderNopCloser(bytes.NewBufferString(f.p.TryWait().String()))), nil
}

func (f *TryWaitFile) Remove() error {
	return rh.ErrPerm
}
