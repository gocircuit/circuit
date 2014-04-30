// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"os"
	"path"
)

// Term
type Term struct {
	slash string
	walk []string // path to this term within term subtree
	dir *dir // open directory of this term
}

func openTerm(slash string) (a *Term, err error) {
	a = &Term{slash: slash}
	if a.dir, err = openDir(slash); err != nil {
		return nil, err
	}
	return
}

// Path returns the path of this term within the local file system.
func (a *Term) Path() string {
	return path.Join(append([]string{a.slash}, a.walk...)...)
}

// UseTerm
func (a *Term) Term(walk ...string) (sub *Term) {
	if len(walk) == 0 {
		return a
	}
	switch walk[0] {
	case "chan", "proc", "help":
		panic("subterms not allowed in element directories")
	}
	os.MkdirAll(path.Join(a.Path(), walk[0]), 0777) // TODO: unused directories are gc'd by the circuit daemon
	sub = &Term{
		slash: a.slash,
		walk: append(a.walk, walk[0]),
	}
	var err error
	if sub.dir, err = openDir(sub.Path()); err != nil {
		panic(err)
	}
	return sub.Term(walk[1:]...)
}

// UseChan
func (a *Term) Chan(name string) *Chan {
	local := path.Join(a.Path(), "chan", name)
	os.MkdirAll(local, 0777)
	return openChan(local)
}

// UseProc
func (a *Term) Proc(name string) *Proc {
	local := path.Join(a.Path(), "proc", name)
	os.MkdirAll(local, 0777)
	return openProc(local)
}
