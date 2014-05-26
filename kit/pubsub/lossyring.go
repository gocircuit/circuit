// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package pubsub

import (
	"sync"
)

// LossyRing
type LossyRing struct {
	sync.Mutex
	head int
	tail int
	cyc  []interface{}
}

// Loss is a special value used in place of a contiguous set of lost messages
type Loss struct {
	Count int
}

// MakeLossyRing
func MakeLossyRing(capacity int) *LossyRing {
	return (&LossyRing{}).Clear(capacity)
}

// Clear
func (s *LossyRing) Clear(capacity int) *LossyRing {
	if capacity < 3 {
		panic("too small")
	}
	s.head, s.tail = 0, 0
	s.cyc = make([]interface{}, capacity)
	return s
}

func (s *LossyRing) Len() int {
	s.Lock()
	defer s.Unlock()
	return s.tail - s.head
}

// Send returns true if and only if the message was stored in the ring.
func (s *LossyRing) Send(v interface{}) (noloss bool) {
	s.Lock()
	defer s.Unlock()
	noloss = true
	if s.head-s.tail == len(s.cyc) {
		if loss, ok := s.cyc[s.tail%len(s.cyc)].(Loss); ok {
			s.cyc[(s.tail+1)%len(s.cyc)] = Loss{loss.Count + 1}
		} else {
			s.cyc[(s.tail+1)%len(s.cyc)] = Loss{2}
		}
		s.tail++
		noloss = false
	}
	s.cyc[s.head%len(s.cyc)] = v
	s.head++
	return
}

func (s *LossyRing) Recv() (v interface{}, ok bool) {
	s.Lock()
	defer s.Unlock()
	if s.tail == s.head {
		return nil, false
	}
	s.tail++
	return s.cyc[(s.tail-1)%len(s.cyc)], true
}
