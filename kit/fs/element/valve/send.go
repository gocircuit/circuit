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
	"github.com/gocircuit/circuit/kit/interruptible"
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
	if flag.Attr != rh.WriteOnly && !flag.Truncate {
		return nil, rh.ErrPerm
	}
	r, w := interruptible.Pipe() // Send –> Recv
	if err := f.v.Send(r, intr); err != nil {
		w.Close()
		r.Close()
		return nil, RhError(err)
	}
	return file.NewOpenInterruptibleWriterFile(w), nil
}

func (f *SendFile) Remove() error {
	return rh.ErrPerm
}
