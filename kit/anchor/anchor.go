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

type Anchor interface {
	Walk([]string) Anchor
	Attach(string, interface{})
	Detach(string)
}

type anchor struct {
	parent *anchor
	name string
	sync.Mutex
	children map[string]interface{}
	nhandle int
}

func newAnchor(parent *anchor, name string) *anchor {
	return &anchor{
		parent: parent,
		name: name,
		children: make(map[string]interface{}),
	}
}

func (a *anchor) Use() *handle {
	a.Lock()
	defer a.Unlock()
	a.nhandle++
	return newHandle(a)
}

func (a *anchor) recycle() {
	a.Lock()
	defer a.Unlock()
	a.nhandle--
	if a.unnecessary() && a.parent != nil {
		a.parent.Detach(a.name)
		a.parent = nil // to catch bugs
	}
}

func (a *anchor) unnecessary() bool {
	return a.nhandle == 0 && len(a.children) == 0
}

func (a *anchor) Walk(walk []string) *anchor {
	if len(walk) == 0 {
		return a
	}
	a.Lock()
	defer a.Unlock()
	q, ok := a.children[walk[0]]
	if !ok {
		q = newAnchor(a, walk[0])
		a.children[walk[0]] = q
	}
	u, ok := q.(*anchor)
	if !ok {
		return nil
	}
	return u.Walk(walk[1:])
}

func (a *anchor) Attach(name string, v interface{}) {
	a.Lock()
	defer a.Unlock()
	a.children[name] = v
}

func (a *anchor) Detach(name string) {
	a.Lock()
	defer a.Unlock()
	delete(a.children, name)
}

func (a *anchor) Children() map[string]interface{} {
	a.Lock()
	defer a.Unlock()
	r := make(map[string]interface{})
	for k, v := range a.children {
		r[k] = v
	}
	return r
}
