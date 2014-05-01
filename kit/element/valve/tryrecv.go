// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

type TryRecvFile struct {
	v *Valve
}

func NewTryRecvFile(v *Valve) file.File {
	return &TryRecvFile{v: v}
}

func (f *TryRecvFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *TryRecvFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.ReadOnly {
		return nil, rh.ErrPerm
	}
	r, err := f.v.TryRecv()
	if err != nil {
		return nil, err
	}
	return file.NewOpenInterruptibleReaderFile(r), nil
}

func (f *TryRecvFile) Remove() error {
	return rh.ErrPerm
}
