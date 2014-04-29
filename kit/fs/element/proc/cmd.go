// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"bytes"
	"encoding/json"

	"github.com/gocircuit/circuit/kit/iomisc"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

// Cmd …
type Cmd struct {
	Env  []string `json:"env"`
	Path string   `json:"path"`
	Args []string `json:"args"`
}

func ParseCmd(src string) (*Cmd, error) {
	x := &Cmd{}
	if err := json.Unmarshal([]byte(src), x); err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Cmd) String() string {
	b, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		panic(0)
	}
	return string(b)
}

// CmdFile …
type CmdFile struct {
	p *Proc
}

func NewCmdFile(p *Proc) file.File {
	return &CmdFile{p: p}
}

func (f *CmdFile) Perm() rh.Perm {
	return 0666 // rw-rw-rw-
}

func (f *CmdFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Truncate {
		return rh.NopClunkerFID{}, nil
	}
	switch flag.Attr {
	case rh.ReadOnly:
		return file.NewOpenReaderFile(
			iomisc.ReaderNopCloser(bytes.NewBufferString(f.p.GetCmd().String())),
		), nil
	case rh.WriteOnly:
		return file.NewOpenWriterFile(&runWriteFile{p: f.p}), nil
	}
	return nil, rh.ErrPerm
}

func (f *CmdFile) Remove() error {
	return rh.ErrPerm
}

type runWriteFile struct {
	p *Proc
	bytes.Buffer
}

func (w *runWriteFile) Close() error {
	w.p.ErrorFile.Clear()
	i, err := ParseCmd(w.Buffer.String())
	if err != nil {
		w.p.ErrorFile.Set("execution description not recognized as JSON")
		return rh.ErrClash
	}
	w.p.SetCmd(i)
	if err := w.p.Start(); err != nil {
		return rh.ErrClash
	}
	return nil
}
