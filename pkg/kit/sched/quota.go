// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package sched

import (
	"io"
	"sync"
)

// Quota behaves like sync.Mutex, except it allows up to a given limit of concurrent lock holders.
type Quota struct {
	lk     sync.Mutex
	ch     chan struct{}
	closed bool
}

// NewQuota creates a new quota with limit m.
func NewQuota(m int) *Quota {
	return (&Quota{}).Init(m)
}

// Init resets this limiter and sets its limit to m.
func (q *Quota) Init(m int) *Quota {
	q.ch = make(chan struct{}, m)
	for i := 0; i < m; i++ {
		q.ch <- struct{}{}
	}
	return q
}

// Lock blocks until there are fewer than limit unclosed sessions.
// A session begins when Open returns.
func (q *Quota) Begin() error {
	// Take an execution permit
	if _, ok := <-q.ch; !ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// End undoes one previous Begin.
func (q *Quota) End() {
	q.lk.Lock()
	defer q.lk.Unlock()
	// Replace the execution permit
	if q.closed {
		return
	}
	q.ch <- struct{}{}
}

func (q *Quota) Close() {
	q.lk.Lock()
	defer q.lk.Unlock()
	if q.closed {
		return
	}
	close(q.ch)
	q.closed = true
}
