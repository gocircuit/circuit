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

type StdinFile struct {
	p *Proc
}

func NewStdinFile(p *Proc) file.File {
	return &StdinFile{p: p}
}

func (f *StdinFile) Perm() rh.Perm {
	return 0222 // -w--w--w-
}

func (f *StdinFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.WriteOnly /* && !flag.Truncate */ {
		return nil, rh.ErrPerm
	}
	return file.NewOpenInterruptibleWriterFile(f.p.Stdin), nil
}

func (f *StdinFile) Remove() error {
	return rh.ErrPerm
}
