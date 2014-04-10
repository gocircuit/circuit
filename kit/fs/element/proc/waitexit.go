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

type WaitFile struct {
	p *Proc
}

func NewWaitFile(p *Proc) file.File {
	return &WaitFile{p: p}
}

func (f *WaitFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *WaitFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.ReadOnly {
		return nil, rh.ErrPerm
	}
	err := f.p.Wait(intr)
	if err == rh.ErrIntr {
		return nil, err
	}
	var msg string
	if err != nil {
		msg = err.Error()
	}
	return file.NewOpenReaderFile(iomisc.ReaderNopCloser(bytes.NewBufferString(msg))), nil
}

func (f *WaitFile) Remove() error {
	return rh.ErrPerm
}
