// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

const selDir = "select"

type sel struct {
	namespace *Namespace
	name      string
	dir       *Dir
	clause    []Clause
}

func makeSel(namespace *Namespace, name string, clause []Clause) (s *sel, err error) {
	s = &sel{
		namespace: namespace,
		name:      name,
		clause:    clause,
	}
	if err = os.Mkdir(s.Path(), 0777); err != nil {
		return nil, err
	}
	if err = ioutil.WriteFile(path.Join(s.Path(), "select"), marshalClauses(clause), 0222); err != nil {
		os.Remove(s.Path())
		return nil, err
	}
	if s.dir, err = OpenDir(s.Path()); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *sel) Path() string {
	return path.Join(s.namespace.Path(), selDir, s.name)
}

func (s *sel) Wait() (branch int, value interface{}) {
	??
}

func marshalClauses(clauses []Clause) []byte {
	var w bytes.Buffer
	for i, c := range clauses {
		switch t := c.(type) {
		}
	}
	return w.Bytes()
}