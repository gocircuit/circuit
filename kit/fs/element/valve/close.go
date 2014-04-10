// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"bytes"
	"strings"

	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

type CloseFile struct {
	v *Valve
}

func NewCloseFile(v *Valve) file.File {
	return &CloseFile{v: v}
}

func (f *CloseFile) Perm() rh.Perm {
	return 0222 // w--w--w--
}

func (f *CloseFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.WriteOnly {
		return nil, rh.ErrPerm
	}
	return file.NewOpenWriterFile(&closeFile{v: f.v}), nil
}

func (f *CloseFile) Remove() error {
	return rh.ErrPerm
}

// closeFile is a io.WriteCloser which closes the valve on Close, if "close" was written to it.
type closeFile struct {
	v *Valve
	bytes.Buffer
}

func (f *closeFile) Close() error {
	f.v.ErrorFile.Clear()
	if strings.TrimSpace(f.Buffer.String()) == "close" {
		f.v.Close()
		return nil
	}
	f.v.ErrorFile.Set("data written to the close file is not \"close\"")
	return rh.ErrClash
}
