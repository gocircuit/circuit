// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package sel

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

type AbortFile struct {
	s *Select
}

func NewAbortFile(s *Select) file.File {
	return &AbortFile{s: s}
}

func (f *AbortFile) Perm() rh.Perm {
	return 0222 // w--w--w--
}

func (f *AbortFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.WriteOnly {
		return nil, rh.ErrPerm
	}
	return file.NewOpenWriterFile(&abortFile{s: f.s}), nil
}

func (f *AbortFile) Remove() error {
	return rh.ErrPerm
}

// abortFile is a io.WriteCloser which aborts the valve on Abort, if "abort" was written to it.
type abortFile struct {
	s *Select
	bytes.Buffer
}

func (f *abortFile) Close() error {
	f.s.ErrorFile.Clear()
	var cmd string
	if err := json.Unmarshal(f.Buffer.Bytes(), &cmd); err != nil {
		f.s.ErrorFile.Set("cannot recognize JSON")
		return rh.ErrClash
	}
	if strings.TrimSpace(cmd) == "abort" {
		return f.s.Abort()
	}
	f.s.ErrorFile.Set("command given is not “abort”")
	return rh.ErrClash
}
