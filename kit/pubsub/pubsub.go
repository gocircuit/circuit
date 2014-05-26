// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package pubsub

import (
	"container/list"
	"runtime"
	"sync"
	
	"github.com/gocircuit/circuit/use/circuit"
)

// PubSub…
type PubSub struct {
	sync.Mutex
	member map[int]*queue
	n int
}

func NewPubSub(src <-chan interface{}) *PubSub {
	ps := &PubSub{
		member: make(map[int]*queue),
	}
	go ps.loop(src)
	return ps
}

func (ps *PubSub) loop(src <-chan interface{}) {
	for {
		v, ok := <-src
		if !ok {
			ps.Lock()
			for _, q := range ps.member {
				q.clunk()
			}
			ps.Unlock()
			return
		}
		ps.Lock()
		for _, q := range ps.member {
			q.distribute(v)
		}
		ps.Unlock()
	}
}

func (ps *PubSub) Subscribe() *Subscription {
	ps.Lock()
	defer ps.Unlock()
	q := newQueue(ps, ps.n)
	ps.member[q.id] = q
	ps.n++
	return q.use()
}

func (ps *PubSub) scrub(id int) {
	ps.Lock()
	defer ps.Unlock()
	s, ok := ps.member[id]
	if !ok || s.isBusy() {
		return
	}
	delete(ps.member, id)
}

// queue…
type queue struct {
	ps *PubSub
	id int
	ch1 chan<- interface{} // disribute() => loop()
	ch2 <-chan interface{} // loop() => consume()
	sync.Mutex
	nref int // number of references to this queue
	pend int // number of buffered messages
	closed bool // true if the source channel has reached EOF
}

func newQueue(ps *PubSub, id int) *queue {
	ch1 := make(chan interface{}, 1)
	ch2 := make(chan interface{}, 1)
	q := &queue{
		ps: ps, 
		id: id, 
		ch1: ch1, 
		ch2: ch2,
	}
	go q.loop(ch1, ch2)
	return q
}

func (q *queue) addPend(d int) {
	q.Lock()
	defer q.Unlock()
	q.pend += d
}

type Stat struct {
	Pending int
	Closed bool
}

func (q *queue) Stat() Stat {
	q.Lock()
	defer q.Unlock()
	return Stat{
		Pending: q.pend,
		Closed: q.closed,
	}
}

func (q *queue) clunk() {
	close(q.ch1)
}

func (q *queue) distribute(v interface{}) {
	q.ch1 <- v
}

func (q *queue) loop(ch1 <-chan interface{}, ch2 chan<- interface{}) {
	var l list.List
__F1:
	for {
		if w := l.Back(); w != nil {
			select {
			case v, ok := <-ch1: // distribute
				if !ok {
					q.closed = true
					break __F1
				}
				l.PushFront(v)
				q.addPend(1)
			case ch2 <- w.Value: // consume
				l.Remove(w)
				q.addPend(-1)
			}
		} else {
			v, ok := <- ch1
			if !ok {
				break __F1
			}
			l.PushFront(v)
			q.addPend(1)
		}
	}
	// After ch1 has been closed
	for {
		w := l.Back()
		if w == nil {
			close(ch2)
			return
		}
		ch2 <- w.Value
		l.Remove(w)
		q.addPend(-1)
	}
}

func (q *queue) isBusy() bool {
	q.Lock()
	defer q.Unlock()
	return q.nref != 0
}

func (q *queue) recycle() {
	q.Lock()
	defer q.Unlock()
	q.nref--
	if q.nref != 0 {
		return
	}
	go q.ps.scrub(q.id)
}

func (q *queue) use() *Subscription {
	q.Lock()
	defer q.Unlock()
	q.nref++
	s := &Subscription{q}
	runtime.SetFinalizer(s, func(s2 *Subscription) {
		q.recycle()
	})
	return s
}

func (q *queue) Consume() (v interface{}, ok bool) {
	v, ok = <-q.ch2
	return
}

// Subscription is the user's interface to consuming messages from a topic.
type Subscription struct {
	*queue // Consume(), Stat()
}

func init() {
	circuit.RegisterValue(&Subscription{})
}

func (s *Subscription) X() circuit.X {
	return circuit.Ref(s)
}

func (s *Subscription) Stat() Stat {
	return s.queue.Stat()
}

func (s *Subscription) Consume() (interface{}, bool) {
	return s.queue.Consume()
}
