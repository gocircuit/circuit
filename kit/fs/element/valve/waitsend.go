// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"bytes"

	"github.com/gocircuit/circuit/kit/iomisc"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

type WaitSendFile struct {
	v *Valve
}

func NewWaitSendFile(v *Valve) file.File {
	return &WaitSendFile{v: v}
}

func (f *WaitSendFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *WaitSendFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.ReadOnly {
		return nil, rh.ErrPerm
	}
	if err := f.v.WaitSend(intr); err != nil {
		return nil, err
	}
	return file.NewOpenReaderFile(iomisc.ReaderNopCloser(bytes.NewBufferString(""))), nil
}

func (f *WaitSendFile) Remove() error {
	return rh.ErrPerm
}
