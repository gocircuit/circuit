// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tube

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
	return &LossyRing{
		head: 0,
		tail: 0,
		cyc:  make([]interface{}, capacity),
	}
}

// Clear
func (s *LossyRing) Clear(length, capacity int) *LossyRing {
	s.head, s.tail = 0, 0
	s.cyc = make([]interface{}, capacity)
	return s
}

func (s *LossyRing) Len() int {
	s.Lock()
	defer s.Unlock()
	return s.tail - s.head
}

func (s *LossyRing) Send(v interface{}) {
	s.Lock()
	defer s.Unlock()
	switch {
	case s.head-s.tail < len(s.cyc):
		s.cyc[s.head%len(s.cyc)] = v
		s.head++
		return
	case s.head-s.tail == len(s.cyc):
		if loss, ok := s.cyc[s.tail%len(s.cyc)].(Loss); ok { // The next message to be received is a loss
			s.cyc[(s.tail+1)%len(s.cyc)] = Loss{loss.Count + 1}
		} else { // The next message to be received is not a loss
			s.cyc[(s.tail+1)%len(s.cyc)] = Loss{2}
		}
		s.cyc[s.tail%len(s.cyc)] = v
		s.head++
		s.tail++
		return
	}
	panic("x")
}

func (s *LossyRing) Recv() (interface{}, bool) {
	s.Lock()
	defer s.Unlock()
	if s.tail > s.head {
		s.tail++
		return s.cyc[(s.tail-1)%len(s.cyc)], true
	}
	return nil, false
}
