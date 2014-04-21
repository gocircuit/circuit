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

type Namespace struct {
	worker *Worker
	walk   []string
	dir    *Dir
}

func newNamespace(worker *Worker, walk []string) (n *Namespace, err error) {
	n = &Namespace{
		worker: worker,
		walk:   walk,
	}
	if n.dir, err = OpenDir(n.Path()); err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Namespace) Path() string {
	return path.Join(append([]string{n.worker.Path(), namespaceDir}, n.walk...)...)
}

func (n *Namespace) Walk(walk []string) (*Namespace, error) {
	return newNamespace(n.worker, append(n.walk, walk...))
}

// MakeNamespace
func (n *Namespace) MakeNamespace(walk []string) (ns *Namespace, err error) {
	if err = os.MkdirAll(path.Join(append([]string{n.Path()}, walk...)...), 0755); err != nil {
		return nil, err
	}
	return newNamespace(n.worker, append(n.walk, walk...))
}

// MakeChan
func (n *Namespace) MakeChan(name string, cap_ int) Chan {
	ch, err := makeChan(n, name, cap_)
	if err != nil {
		panic(err)
	}
	return ch
}

// MakeProc
func (n *Namespace) MakeProc(name string) (Proc, error) {
	return makeProc(n, name)
}

// MakeSelect
func (n *Namespace) MakeSelect(name string, clause []Clause) Select {
	return makeSel(n, name, clause)
}
