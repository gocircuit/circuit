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

	"github.com/gocircuit/circuit/kit/iomisc"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

type SelectFile struct {
	s *Select
}

func NewSelectFile(s *Select) file.File {
	return &SelectFile{s: s}
}

func (f *SelectFile) Perm() rh.Perm {
	return 0666 // rw-rw-rw-
}

func marshal(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (f *SelectFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	switch flag.Attr {
	case rh.ReadOnly:
		return file.NewOpenReaderFile(
			iomisc.ReaderNopCloser(bytes.NewBufferString(marshal(f.s.GetClauses()))),
		), nil
	case rh.WriteOnly:
		return file.NewOpenWriterFile(&runWriteFile{s: f.s}), nil
	}
	return nil, rh.ErrPerm
}

func (f *SelectFile) Remove() error {
	return rh.ErrPerm
}

type runWriteFile struct {
	s *Select
	bytes.Buffer
}

func (w *runWriteFile) Close() (err error) {
	w.s.ErrorFile.Clear()
	println("parsing json", w.Buffer.String())
	var clauses []Clause
	if err = json.Unmarshal(w.Buffer.Bytes(), &clauses); err != nil {
		w.s.ErrorFile.Set("cannot recognize JSON clauses structure")
		return rh.ErrClash
	}
	if err = w.s.Select(clauses); err != nil {
		return rh.ErrClash
	}
	return nil
}

type Clause struct {
	Op   string `json:"op"`
	File string `json:"file"`
}
