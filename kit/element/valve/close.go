// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"bytes"
	"encoding/json"
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
	var cmd string
	if err := json.Unmarshal(f.Buffer.Bytes(), &cmd); err != nil {
		f.v.ErrorFile.Set("cannot recognize JSON")
		return rh.ErrClash
	}
	if strings.TrimSpace(cmd) == "close" {
		return f.v.Close()
	}
	f.v.ErrorFile.Set("command given is not “close”")
	return rh.ErrClash
}
