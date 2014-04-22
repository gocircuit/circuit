// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/gocircuit/circuit/kit/fs/element/sel"
)

const selDir = "select"

type sel struct {
	namespace *Namespace
	name      string
	dir       *Dir
	sync.Mutex
	clause    []Clause
}

type unblock struct {
	Branch int
	Value  interface{}
	Error  error
}

func makeSel(namespace *Namespace, name string) (s *sel, err error) {
	s = &sel{
		namespace: namespace,
		name:      name,
	}
	if err = os.Mkdir(s.Path(), 0777); err != nil {
		return nil, err
	}
	if s.dir, err = OpenDir(s.Path()); err != nil {
		os.Remove(s.Path())
		return nil, err
	}
	return c, nil
}

// Path returns the path of this select element in the local circuit file system.
func (s *sel) Path() string {
	return path.Join(s.namespace.Path(), selDir, s.name)
}

// ??
func (s *sel) Start(clause []Clause) error {
	if err = ioutil.WriteFile(path.Join(s.Path(), "select"), encodeClauses(clause), 0222); err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	if s.clause != nil {
		panic("selection already started")
	}
	s.clause = clause
}

func encodeClauses(clauses []Clause) []byte {
	var waitfiles []string
	for i, c := range clauses {
		switch t := c.(type) {
		case ClauseSend:
			waitfiles = append(waitfiles, path.Join(t.Chan.Path(), "waitsend")
		case ClauseReceive:
			waitfiles = append(waitfiles, path.Join(t.Chan.Path(), "waitrecv")
		case ClauseExit:
			waitfiles = append(waitfiles, path.Join(t.Proc.Path(), "waitexit")
		}
	}
	buf, err := json.Marshal(waitfiles)
	if err != nil {
		panic(0)
	}
	return buf
}

// hasDefault returns true if one of clause is ClauseDefault.
func hasDefault(clause []Clause) bool {
	for _, c := range clause {
		if _, ok := c.(ClauseDefault); ok {
			return true
		}
	}
	return false
}

var ExitOK = errors.New("process exit ok")
var ExitError = errors.New("process exit error")
var Closed = errors.New("channel closed")

// Wait blocks until one of the select branches unblocks.
// branch indicates the index of the clause that unblocked.
// value is an io.WriteCloser if the unblocking clause is a channel send;
// value is an io.ReadCloser if unblocking clause is a channel receive and the channel is still open;
// value equals Closed if the unblocking clause is a channel receive and the channel is closed;
// value equals ExitOK or ExitError if the unblocking clause is a successful or erroneous process exit.
// value can be nil in the case of channel receive and channel send, if a racing select or channel operation
// snatches the unblocked channel operation first; 
// TODO(petar): the latter condition could be wrapped inside a spin loop
func (s *sel) Wait() (branch int, value interface{}) {
	b, err := ioutil.ReadFile(path.Join(s.Path(), "wait"))
	if os.IsNotExist(err) { // a missing file indicates a dead circuit worker; we panic for those by convention
		panic(err)
	}
	var r sel.Result
	if err = json.Unmarshal(b, &r); err != nil {
		panic(err)
	}
	if r.Branch < 0 { // happens if select is interrupted by removal
		return -1, nil
	}
	//
	s.Lock()
	defer s.Unlock()
	//
	clause := r.clause[r.Branch]
	switch t := clause.(type) {
	case ClauseSend:
		??
	case ClauseReceive:
		??
	case ClauseExit:
		??
	}
	panic(0)
}

// Remove removes the select element from the circuit environment (and local file system).
func (s *sel) Remove() error {
	return os.Remove(s.Path())
}
