// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"sync"
)

type anchor struct {
	parent *anchor
	name string
	sync.Mutex
	children map[string]*anchor
	nhandle int
	elem Element // chan, proc, etc
}

type Element interface {
	IsDone() bool
	Scrub()
}

func newAnchor(parent *anchor, name string) *anchor {
	return &anchor{
		parent: parent,
		name: name,
		children: make(map[string]*anchor),
	}
}

func (a *anchor) scrub(name string) {
	a.Lock()
	defer a.Unlock()
	delete(a.children, name)
}

func (a *anchor) recycle() {
	a.Lock()
	defer a.Unlock()
	a.nhandle--
	if a.unnecessary() && a.parent != nil {
		a.parent.scrub(a.name)
		a.parent = nil // catch bugs
	}
}

func (a *anchor) unnecessary() bool {
	return a.nhandle == 0 && (a.elem == nil || a.elem.IsDone())
}

func (a *anchor) Walk(walk []string) Anchor {
	a.Lock()
	defer a.Unlock()
	if len(walk) == 0 {
		a.nhandle++
		return newHandle(a)
	}
	q, ok := a.children[walk[0]]
	if !ok {
		q = newAnchor(a, walk[0])
		a.children[walk[0]] = q
	}
	return q.Walk(walk[1:])
}

func (a *anchor) Content() (Element, map[string]Anchor) {
	a.Lock()
	defer a.Unlock()
	r := make(map[string]Anchor)
	for k, v := range a.children {
		r[k] = v.Walk(nil)
	}
	return a.elem, r
}

func (a *anchor) Set(elem Element) {
	a.Lock()
	defer a.Unlock()
	a.elem = elem
}
