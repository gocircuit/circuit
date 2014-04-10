// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chain

import (
	"bytes"
	"fmt"
	"net"
	gosync "sync"

	cirsync "github.com/gocircuit/circuit/kit/sync"
)

type Stat struct {
	p cirsync.Publisher
	gosync.Mutex
	dial, accept map[chainID]*ConnStat
	summary      Summary
}

type ConnStat struct {
	Addr net.Addr
}

type Summary struct {
	NDial   int
	NAccept int
}

func (s *Stat) Init() {
	s.dial = make(map[chainID]*ConnStat)
	s.accept = make(map[chainID]*ConnStat)
}

func (s *Stat) Dump() string {
	s.Lock()
	defer s.Unlock()
	//
	var w bytes.Buffer
	for k, v := range s.dial {
		fmt.Fprintf(&w, "D—%s—>%s\n", k, v.Addr)
	}
	for k, v := range s.accept {
		fmt.Fprintf(&w, "A—%s<—%s\n", k, v.Addr)
	}
	return w.String()
}

// Dial

func (s *Stat) addDC(id chainID, addr net.Addr) {
	s.Lock()
	defer s.Unlock()
	defer s.publish() // Publish before unlocking to get consistent representation of the events
	//
	if _, ok := s.dial[id]; ok {
		panic("u")
	}
	s.dial[id] = &ConnStat{Addr: addr}
	//
	s.summary.NDial++
}

func (s *Stat) scrubDC(id chainID) {
	s.Lock()
	defer s.Unlock()
	defer s.publish() // Publish before unlocking to get consistent representation of the events
	//
	if _, ok := s.dial[id]; !ok {
		panic("u")
	}
	delete(s.dial, id)
	//
	s.summary.NDial--
}

// Accept

func (s *Stat) addAC(id chainID, addr net.Addr) {
	s.Lock()
	defer s.Unlock()
	defer s.publish() // Publish before unlocking to get consistent representation of the events
	//
	if _, ok := s.accept[id]; ok {
		panic("u")
	}
	s.accept[id] = &ConnStat{Addr: addr}
	//
	s.summary.NAccept++
}

func (s *Stat) scrubAC(id chainID) {
	s.Lock()
	defer s.Unlock()
	defer s.publish() // Publish before unlocking to get consistent representation of the events
	//
	if _, ok := s.accept[id]; !ok {
		panic("u")
	}
	delete(s.accept, id)
	//
	s.summary.NAccept--
}

// Publishing
func (s *Stat) publish() {
	var u = s.summary
	s.p.Publish(&u)
}

func (s *Stat) Subscribe() *Subscriber {
	return (*Subscriber)(s.p.Subscribe())
}

// Subscriber
type Subscriber cirsync.Subscriber

func (s *Subscriber) Wait() *Summary {
	v := (*cirsync.Subscriber)(s).Wait()
	if v == nil {
		return nil
	}
	return v.(*Summary)
}
