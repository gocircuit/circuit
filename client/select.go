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
	"math/rand"
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
	result    <-chan *unblock
}

type unblock struct {
	Branch int
	Value  interface{}
	Error  error
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

// Wait blocks until one of the select branches unblocks.
// branch indicates the index of the clause that unblocked.
// value is an io.WriteCloser if the unblocking clause is a channel send;
// value is an io.ReadCloser if unblocking clause is a channel receive and the channel is still open;
// value equals nil if the unblocking clause is a channel receive and the channel is closed;
// value is nil if the unblocking clause is a successful process exit, otherwise it is a non-nil error.
func (s *sel) Wait() (branch int, value interface{}) {
	r, ok := <-s.result
	if !ok {
		return -1, nil // all clauses blocked and selection has default clause
	}
	if r.Error != nil { // the worker directory hosting this selection is gone (worker is dead)
		panic(r.Error)
	}
	return r.Branch, r.Value
}

// start initiates the selection
func (s *sel) start(clause []Clause) {
	if hasDefault(clause) {
		s.startNonBlocking(clause)
	} else {
		s.startBlocking(clause)
	}
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

func (s *sel) startNonBlocking(clause []Clause) {
	var i int
	perm := rand.Perm(len(clause))
	//
	defer close(s.result)
	defer func() {
		if r := recover(); r != nil {
			s.result <- &unblock{
				Branch: perm[i],
				Error:  r,
			}
		}
	}()
	//
	for i, _ = range clause {
		c := clause[perm[i]]
		switch t := c.(type) {
		case ClauseDefault:
			// skip

		case ClauseSend:
			s.result <- &unblock{
				Branch: perm[i],
				Value:  t.Chan.TrySend(),
			}
			return

		case ClauseReceive:
			result <- &unblock{
				Branch: perm[i],
				Value:  t.Chan.TryRecv(),
			}
			return

		case ClauseExit:
			result <- &unblock{
				Branch: perm[i],
				Value:  t.Proc.Stat(),
			}
			return

		}
	}
	// If all clauses are blocking, just close s.result within the defer above.
}

func (s *sel) startBlocking(clause []Clause) {
	perm := rand.Perm(len(clause))
	backchan := make(chan *unblock, len(clause))
	for i, _ = range clause { // fire off waiters
		b := perm[i]
		c := clause[b]
		switch t := c.(type) {
		case ClauseDefault:
			panic(0)
		case ClauseSend:
			go waitSend(backchan, b, t.Chan)
		case ClauseReceive:
			go waitRecv(backchan, b, t.Chan)
		case ClauseExit:
			go waitExit(backchan, b, t.Proc)
		}
	}
	for { // spin after wait until operation completes
		??
	}
}

func waitSend(report chan<- *unblock, branch int, ch Chan) {
	defer func() {
		if r := recover(); r != nil {
			select {
			case s.result <- &unblock{
					Branch: branch,
					Error:  r,
				}
			}
		}
	}()
	ch.WaitSend()
	select {
	case ch <- &unblock{Branch: branch}
	}
}

func waitRecv(report chan<- *unblock, branch int, ch Chan) {
	defer func() {
		if r := recover(); r != nil {
			select {
			case s.result <- &unblock{
					Branch: branch,
					Error:  r,
				}
			}
		}
	}()
	ch.WaitRecv()
	select {
	case ch <- &unblock{Branch: branch}
	}
}

func marshalClauses(clauses []Clause) []byte {
	??
	var w bytes.Buffer
	for i, c := range clauses {
		switch t := c.(type) {
		}
	}
	return w.Bytes()
}
