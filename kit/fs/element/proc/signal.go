// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"bytes"
	"strings"

	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

type SignalFile struct {
	p *Proc
}

func NewSignalFile(p *Proc) file.File {
	return &SignalFile{p: p}
}

func (f *SignalFile) Perm() rh.Perm {
	return 0222 // w--w--w--
}

func (f *SignalFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Truncate {
		return rh.NopClunkerFID{}, nil
	}
	if flag.Attr != rh.WriteOnly {
		return nil, rh.ErrPerm
	}
	return file.NewOpenWriterFile(&signalFile{p: f.p}), nil
}

func (f *SignalFile) Remove() error {
	return rh.ErrPerm
}

type signalFile struct {
	p *Proc
	bytes.Buffer
}

func (f *signalFile) Close() error {
	return f.p.Signal(strings.TrimSpace(f.Buffer.String()))
}
