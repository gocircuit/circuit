// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package sync provides various synchronization primitives
package sync

import (
	"sync"
)

// Trigger provides a mechanism for reconciling competing callers, only one of which should succeed in
// obtaining the lock (and taking some action).
type Trigger struct {
	lk       sync.Mutex
	engaged  bool
	nwaiters int
	ch       chan struct{}
}

// Lock attempts to lock the trigger.
// If the trigger is not currently locked, Lock returns instantaneously with true.
// Otherwise, it blocks until the trigger is unlocked by its holder and returns false,
// WITHOUT locking the trigger on behalf of the caller.
func (t *Trigger) Lock() bool {
	t.lk.Lock()
	if t.ch == nil {
		t.ch = make(chan struct{})
	}
	if t.engaged {
		t.nwaiters++
		t.lk.Unlock()
		<-t.ch
		return false
	}
	t.engaged = true
	t.lk.Unlock()
	return true
}

// Unlock unlocks a locked trigger.
func (t *Trigger) Unlock() {
	t.lk.Lock()
	defer t.lk.Unlock()
	if !t.engaged {
		panic("unlocking a non-engaged trigger")
	}
	for t.nwaiters > 0 {
		t.ch <- struct{}{}
		t.nwaiters--
	}
	t.engaged = false
}
