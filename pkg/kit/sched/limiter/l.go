// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package limiter schedules job execution while maintaining an upper limit on concurrency
package limiter

// TODO: Implement Limiter in terms of Quota. Move to sched package.

import (
	"sync"
)

// Limiter schedules go routines for execution, while ensuring that no more than a
// pre-set limit run at any time.
type Limiter struct {
	ch chan struct{}
	wg sync.WaitGroup
}

// New creates a new limiter with limit m.
func New(m int) *Limiter {
	return (&Limiter{}).Init(m)
}

// Init resets this limiter and sets its limit to m.
func (l *Limiter) Init(m int) *Limiter {
	l.ch = make(chan struct{}, m)
	for i := 0; i < m; i++ {
		l.ch <- struct{}{}
	}
	return l
}

// Open blocks until there are fewer than limit unclosed sessions.
// A session begins when Open returns.
func (l *Limiter) Open() {
	// Take an execution permit
	<-l.ch
	l.wg.Add(1)
}

// Close closes a session.
func (l *Limiter) Close() {
	// Replace the execution permit
	l.ch <- struct{}{}
	l.wg.Done()
}

// Go executes the function f when the goroutine limit allows it.
// Go wraps the execution of f around an Open/Close pair.
func (l *Limiter) Go(f func()) {
	l.Open()
	go func() {
		f()
		l.Close()
	}()
}

// Throttle executes copies of f greedily and continuously, making sure that at no
// time the limit is exceeded.
func (l *Limiter) Throttle(f func()) {
	for {
		l.Go(f)
	}
}

// Wait blocks until all unclosed invokations to Open have been closed.
func (l *Limiter) Wait() {
	l.wg.Wait()
}
