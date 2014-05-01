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

type StdoutFile struct {
	p *Proc
}

func NewStdoutFile(p *Proc) file.File {
	return &StdoutFile{p: p}
}

func (f *StdoutFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *StdoutFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.ReadOnly {
		return nil, rh.ErrPerm
	}
	return file.NewOpenInterruptibleReaderFile(f.p.Stdout), nil
}

func (f *StdoutFile) Remove() error {
	return rh.ErrPerm
}
