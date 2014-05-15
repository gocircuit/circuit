// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"runtime"
	"sync"

	"github.com/gocircuit/circuit/use/circuit"
)

type Anchor struct {
	*anchor
}

type anchor struct {
	parent *anchor
	name string
	lk sync.Mutex
	children map[string]*anchor
	nhandle int
	value interface{}
	tx sync.Mutex
}

func (a *anchor) TxLock() {
	a.tx.Lock()
}

func (a *anchor) TxUnlock() {
	a.tx.Unlock()
}

func (a *anchor) use() *Anchor {
	a.lk.Lock()
	defer a.lk.Unlock()
	u := &Anchor{a}
	a.nhandle++
	runtime.SetFinalizer(u, func(u2 *Anchor) {
		u2.recycle()
	})
	return u
}

func (a *anchor) recycle() {
	println("recycle")
	a.lk.Lock()
	defer a.lk.Unlock()
	a.nhandle--
	if !a.busy() && a.parent != nil {
		go a.parent.scrub(a.name)
	}
}

func (a *anchor) Busy() bool {
	a.lk.Lock()
	defer a.lk.Unlock()
	return a.busy()
}

func (a *anchor) busy() bool {
	return a.nhandle > 0 || a.value != nil || len(a.children) > 0
}

func (a *anchor) scrub(name string) {
	a.lk.Lock()
	defer a.lk.Unlock()
	q, ok := a.children[name]
	if !ok {
		return
	}
	if q.Busy() {
		return
	}
	println("del", name)
	delete(a.children, name)
}

func (a *anchor) Walk(walk []string) *Anchor {
	if len(walk) == 0 {
		return a.use()
	}
	a.lk.Lock()
	defer a.lk.Unlock()
	q, ok := a.children[walk[0]]
	if !ok {
		println("makin", walk[0])
		q = newAnchor(a, walk[0])
		a.children[walk[0]] = q
		q.use() // ensures that if q is not used after Walk returns, it will be scrubbed
	}
	return q.Walk(walk[1:])
}

// NewTerm create the root node of a new anchor file system.
func NewTerm(name string) circuit.PermX {
	t := XTerminal{
		(*Terminal)(newAnchor(nil, name).use()),
	}
	return circuit.PermRef(t)
}

func newAnchor(parent *anchor, name string) *anchor {
	return &anchor{
		parent: parent,
		name: name,
		children: make(map[string]*anchor),
	}
}

func (a *anchor) View() (r map[string]*Anchor) {
	a.lk.Lock()
	defer a.lk.Unlock()
	r = make(map[string]*Anchor)
	for n, m := range a.children {
		r[n] = m.use()
	}
	return r
}

func (a *anchor) Set(v interface{}) {
	a.lk.Lock()
	defer a.lk.Unlock()
	a.value = v
	if !a.busy() && a.parent != nil {
		go a.parent.scrub(a.name)
	}
}

func (a *anchor) Get() interface{} {
	a.lk.Lock()
	defer a.lk.Unlock()
	return a.value
}
