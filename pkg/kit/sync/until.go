// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package sync provides various synchronization primitives
package sync

import (
	"runtime"
	"sync"
)

// WaitUntil notifies subscriber of a sequence of events
type WaitUntil struct {
	sync.Mutex
	waiters []Waiter
}

func (wu *WaitUntil) Broadcast() {
	wu.Lock()
	waiters := wu.waiters
	wu.waiters = nil
	wu.Unlock()
	//
	for _, ch := range waiters {
		close(ch)
	}
}

func (wu *WaitUntil) MakeWaiter() Waiter {
	ch := make(Waiter)
	wu.Lock()
	wu.waiters = append(wu.waiters, ch)
	wu.Unlock()
	return ch
}

// Waiter receivers a continuous stream of event notifications
type Waiter chan struct{}

func (w Waiter) Wait() {
	<-w
}

// Publisher allows multiple subscriptions to a continuous stream of published events.
// Subscribers can cancel their subscription.
// The zero value is ready for use.
type Publisher struct {
	s__ sync.Mutex
	s   []*Subscriber
	p__ sync.Mutex
}

func (p *Publisher) Publish(v interface{}) {
	p.p__.Lock()
	defer p.p__.Unlock()
	//
	p.s__.Lock()
	s := p.s
	p.s__.Unlock()
	//
	for _, t := range s {
		t.ch <- v
	}
}

func (p *Publisher) scrub(s *Subscriber) {
	p.s__.Lock()
	defer p.s__.Unlock()
	for i, r := range p.s {
		if r == s {
			n := len(p.s)
			p.s[i] = p.s[n-1]
			p.s = p.s[:n-1]
			return
		}
	}
}

func (p *Publisher) Subscribe() *Subscriber {
	s := &Subscriber{
		p:  p,
		ch: make(chan interface{}, 5),
	}
	p.s__.Lock()
	p.s = append(p.s, s)
	p.s__.Unlock()
	runtime.SetFinalizer(s, func(x *Subscriber) {
		x.Scrub()
	})
	return s
}

// Subscriber
type Subscriber struct {
	p  *Publisher
	ch chan interface{}
}

// Scrub and Wait must be not be called concurrently.
func (s *Subscriber) Wait() interface{} {
	return <-s.ch
}

// Scrub and Wait must be not be called concurrently.
func (s *Subscriber) Scrub() {
	s.p.scrub(s)
	for _, ok := <-s.ch; ok; _, ok = <-s.ch {
		// Drain channel
	}
}
