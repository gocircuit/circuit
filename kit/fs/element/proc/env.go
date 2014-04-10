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

type EnvFile struct {
	p *Proc
}

func NewEnvFile(p *Proc) file.File {
	return &EnvFile{p: p}
}

func (f *EnvFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *EnvFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.ReadOnly {
		return nil, rh.ErrPerm
	}
	return file.NewOpenReaderFile(iomisc.ReaderNopCloser(bytes.NewBufferString(marshal(f.p.GetEnv())))), nil
}

func (f *EnvFile) Remove() error {
	return rh.ErrPerm
}

func marshal(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(b)
}
