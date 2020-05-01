// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"sync"
)

type srvTabl struct {
	sync.Mutex
	name map[string]interface{}
}

func (t *srvTabl) Init() *srvTabl {
	t.Lock()
	defer t.Unlock()
	t.name = make(map[string]interface{})
	return t
}

func (t *srvTabl) Add(name string, receiver interface{}) {
	t.Lock()
	defer t.Unlock()
	if _, present := t.name[name]; present {
		panic("service already listening")
	}
	x := receiver
	t.name[name] = x
}

func (t *srvTabl) Get(name string) interface{} {
	t.Lock()
	defer t.Unlock()
	return t.name[name]
}
