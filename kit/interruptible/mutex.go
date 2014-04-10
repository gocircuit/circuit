// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package interruptible

import (
	"sync"

	"github.com/gocircuit/circuit/kit/fs/rh" // TODO: backwards dep
)

// Mutex is analgous to sync.Mutex, but the lock operation can be interrupted by the locking user.
type Mutex struct {
	lk   sync.Mutex
	wait <-chan struct{}
}

func (m *Mutex) Lock(intr rh.Intr) *Unlocker {
	//
	m.lk.Lock()
	if m.wait == nil { // initial condition
		u := make(chan struct{})
		m.wait = u
		close(u)
	}
	wait := m.wait
	turn := make(chan struct{})
	m.wait = turn
	m.lk.Unlock()
	//
	select {
	case <-wait:
		return &Unlocker{turn}
	case <-intr:
		close(turn)
		return nil
	}
}

type Unlocker struct {
	turn chan<- struct{}
}

func (u *Unlocker) Unlock() {
	close(u.turn)
}
