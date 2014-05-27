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
	up struct {
		sync.Mutex
		src chan<- interface{}
	}
	down struct {
		sum Summarize
		sync.Mutex
		member map[int]*queue
		n int
	}
}

// Summarize returns a list of items meant to summarize the history of the stream so far
// for subscribers joining now.
type Summarize func() []interface{}

func New(sum Summarize) (ps *PubSub) {
	src := make(chan interface{})
	ps = &PubSub{}
	ps.up.src = src
	ps.down.sum = sum
	ps.down.member = make(map[int]*queue)
	go ps.loop(src)
	return
}

// Publish appends a value onto the infinite update stream.
func (ps *PubSub) Publish(v interface{}) {
	ps.up.Lock()
	defer ps.up.Unlock()
	if ps.up.src == nil {
		panic("publish after close")
	}
	ps.up.src <- v
}

// Close terminates, paradoxically, the infinite update stream.
func (ps *PubSub) Close() {
	ps.up.Lock()
	defer ps.up.Unlock()
	if ps.up.src == nil {
		return
	}
	close(ps.up.src)
	ps.up.src = nil
}

// loop churns messages between the publishing entity, using Publish(), 
// and the multiple registered subscriber entities.
func (ps *PubSub) loop(src <-chan interface{}) {
	for {
		v, ok := <-src
		if !ok {
			ps.clunk()
			return
		}
		ps.distribute(v)
	}
}

func (ps *PubSub) distribute(v interface{}) {
	ps.down.Lock()
	defer ps.down.Unlock()
	for _, q := range ps.down.member {
		q.distribute(v)
	}
}

func (ps *PubSub) clunk() {
	ps.down.Lock()
	defer ps.down.Unlock()
	for _, q := range ps.down.member {
		q.clunk()
	}
}

// Subscribe creates a new subscription object, whose interface embodies reading from an infinite stream.
// New subscription can join at any time. The input stream of each individual subscription is pre-loaded
// with a sequence of values summarizing all past history. Subsequent values come from the pubish stream.
// Subscriptions are abandoned on garbage-collection.
func (ps *PubSub) Subscribe() *Subscription {
	ps.down.Lock()
	defer ps.down.Unlock()
	q := newQueue(ps, ps.down.n)
	ps.down.member[q.id] = q
	ps.down.n++
	// Prefix subscription's input stream with a summary of all history until now
	if ps.down.sum != nil {
		for _, v := range ps.down.sum() {
			q.distribute(v)
		}
	}
	return q.use()
}

// scrub removes a subscription queue from the member table, only if 
// all Subscription handles referring to it have been collected.
func (ps *PubSub) scrub(id int) {
	ps.down.Lock()
	defer ps.down.Unlock()
	q, ok := ps.down.member[id]
	if !ok || q.isBusy() {
		return
	}
	delete(ps.down.member, id)
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

func (q *queue) setClosed(v bool) {
	q.Lock()
	defer q.Unlock()
	q.closed = v
}

type Stat struct {
	Pending int
	Closed bool
}

func (q *queue) Peek() Stat {
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

// loop churns messages from the main loop onto the internal buffer of this subscription,
// and from there out to the consumer, as requested by calls to Consume.
func (q *queue) loop(ch1 <-chan interface{}, ch2 chan<- interface{}) {
	var l list.List
__preclose:
	for {
		if w := l.Back(); w != nil {
			select {
			case v, ok := <-ch1: // distribute
				if !ok {
					q.setClosed(true)
					break __preclose
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
				q.setClosed(true)
				break __preclose
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

// use returns a new subscription, which is effectively a handle for this queue.
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
	*queue // Consume(), Peek()
}

func init() {
	circuit.RegisterValue(&Subscription{})
}

func (s *Subscription) X() circuit.X {
	return circuit.Ref(s)
}

func (s *Subscription) Scrub() {}

func (s *Subscription) Peek() Stat {
	return s.queue.Peek()
}

func (s *Subscription) Consume() (interface{}, bool) {
	return s.queue.Consume()
}

// YSubscription is a client wrapper for cross-interface to *Subscription
type YSubscription struct {
	X circuit.X
}

func (y YSubscription) Peek() Stat {
	return y.X.Call("Peek")[0].(Stat)
}

func (y YSubscription) Consume() (interface{}, bool) {
	r := y.X.Call("Consume")
	return r[0], r[1].(bool)
}

func (y YSubscription) IsDone() bool {
	return true
}

func (y YSubscription) Scrub() {}
