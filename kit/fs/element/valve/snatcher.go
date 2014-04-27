// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"sync"
)

// Snatcher allows only one entity, of a collective or racing entities, to take a planted value.
type Snatcher struct {
	h <-chan struct{}
	sync.Mutex
	w interface{} // who snatched the value
}

func NewSnatcher() *Snatcher {
	h := make(chan struct{}, 1)
	h <- struct{}{}
	close(h)
	return &Snatcher{
		h: h,
	}
}

type SnatchResult int
const (
	FirstSnatch SnatchResult = iota
	RepeatSnatch
	RejectSnatch
)

func (s *Snatcher) Snatch(who interface{}) SnatchResult {
	_, ok := <-s.h
	s.Lock()
	defer s.Unlock()
	if ok {
		s.w = who
		return FirstSnatch
	}
	if s.w == who {
		return RepeatSnatch
	}
	return RejectSnatch
}
