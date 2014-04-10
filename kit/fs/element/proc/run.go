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

type RunFile struct {
	p *Proc
}

func NewRunFile(p *Proc) file.File {
	return &RunFile{p: p}
}

func (f *RunFile) Perm() rh.Perm {
	return 0666 // rw-rw-rw-
}

func (f *RunFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	switch flag.Attr {
	case rh.ReadOnly:
		return file.NewOpenReaderFile(
			iomisc.ReaderNopCloser(bytes.NewBufferString(f.p.GetInit().String())),
		), nil
	case rh.WriteOnly:
		return file.NewOpenWriterFile(&runWriteFile{p: f.p}), nil
	}
	return nil, rh.ErrPerm
}

func (f *RunFile) Remove() error {
	return rh.ErrPerm
}

type runWriteFile struct {
	p *Proc
	bytes.Buffer
}

func (w *runWriteFile) Close() error {
	w.p.ErrorFile.Clear()
	i, err := ParseInit(w.Buffer.String())
	if err != nil {
		w.p.ErrorFile.Set("data written not JSON")
		return rh.ErrClash
	}
	w.p.SetInit(i)
	if err := w.p.Start(); err != nil {
		return rh.ErrClash
	}
	return nil
}
