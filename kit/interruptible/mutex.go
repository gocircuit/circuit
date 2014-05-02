// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package interruptible

import (
	"sync"
)

// Mutex is analgous to sync.Mutex, but the lock operation can be interrupted by the locking user.
type Mutex struct {
	initlk sync.Mutex
	lock   <-chan struct{}
	unlock chan<- struct{}
}

func (m *Mutex) init() {
	m.initlk.Lock()
	defer m.initlk.Unlock()
	if m.lock != nil {
		return
	}
	ch := make(chan struct{}, 1)
	ch <- struct{}{} // first lock lease
	m.lock, m.unlock = ch, ch
}

func (m *Mutex) Lock(intr Intr) *Unlocker {
	m.init()
	select {
	case <-m.lock:
		return &Unlocker{m}
	case <-intr:
		return nil
	}
}

func (m *Mutex) TryLock() *Unlocker {
	m.init()
	select {
	case <-m.lock:
		return &Unlocker{m}
	default:
		return nil
	}
}

type Unlocker struct {
	mutex *Mutex
}

func (u *Unlocker) Unlock() {
	select {
	case u.mutex.unlock <- struct{}{}:
	}
}
