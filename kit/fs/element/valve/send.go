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

type SendFile struct {
	v *Valve
}

func NewSendFile(v *Valve) file.File {
	return &SendFile{v: v}
}

func (f *SendFile) Perm() rh.Perm {
	return 0222 // -w--w--w-
}

func (f *SendFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	//println("open/send", flag.String())
	if flag.Attr != rh.WriteOnly {
		return nil, rh.ErrPerm
	}
	w, err := f.v.Send(intr)
	if err != nil {
		return nil, err
	}
	return file.NewOpenInterruptibleWriterFile(w), nil
}

func (f *SendFile) Remove() error {
	return rh.ErrPerm
}
