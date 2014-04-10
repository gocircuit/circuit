// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

type StderrFile struct {
	p *Proc
}

func NewStderrFile(p *Proc) file.File {
	return &StderrFile{p: p}
}

func (f *StderrFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *StderrFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.ReadOnly {
		return nil, rh.ErrPerm
	}
	return file.NewOpenInterruptibleReaderFile(f.p.Stderr), nil
}

func (f *StderrFile) Remove() error {
	return rh.ErrPerm
}
