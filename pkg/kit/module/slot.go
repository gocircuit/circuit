// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package module provides a mechanism for linking an implementation package to a declaration package
package module

import (
	"sync"
)

// Slot is a synchronized interface value, which can be set once and read many times
type Slot struct {
	Name string
	lk   sync.Mutex
	v    interface{}
}

// Set sets the value to v
func (j *Slot) Set(v interface{}) {
	j.lk.Lock()
	defer j.lk.Unlock()
	if j.v != nil {
		panic(j.Name + " already set")
	}
	j.v = v
}

// Get returns this value
func (j *Slot) Get() interface{} {
	j.lk.Lock()
	defer j.lk.Unlock()
	if j.v == nil {
		panic(j.Name + " not set")
	}
	return j.v
}
