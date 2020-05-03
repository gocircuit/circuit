// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tissue

import (
	"sync"

	"github.com/hoijui/circuit/pkg/kit/lang"
)

// Neighborhood is a set of perm cross-interfaces.
type Neighborhood struct {
	sync.Mutex
	open map[interface{}]Avatar
}

// NewNeighborhood creates a new rotor.
func NewNeighborhood() *Neighborhood {
	return &Neighborhood{
		open: make(map[interface{}]Avatar),
	}
}

func (nh *Neighborhood) Add(av Avatar) {
	nh.Lock()
	defer nh.Unlock()
	nh.open[av.ID] = av
}

func (nh *Neighborhood) Scrub(key lang.ReceiverID) (Avatar, bool) {
	nh.Lock()
	defer nh.Unlock()
	av, ok := nh.open[key]
	delete(nh.open, key)
	return av, ok
}

func (nh *Neighborhood) ScrubRandom() (Avatar, bool) {
	nh.Lock()
	defer nh.Unlock()
	for key, av := range nh.open {
		delete(nh.open, key)
		return av, true
	}
	return Avatar{}, false
}

// View returns a list of all Avatars in the rotor.
func (nh *Neighborhood) View() []Avatar {
	nh.Lock()
	defer nh.Unlock()
	open := make([]Avatar, 0, len(nh.open))
	for _, av := range nh.open {
		open = append(open, av)
	}
	return open
}

// Len returns the number of Avatars in the rotor.
func (nh *Neighborhood) Len() int {
	nh.Lock()
	defer nh.Unlock()
	return len(nh.open)
}

// Choose returns a randomly chosen Avatar.
func (nh *Neighborhood) Choose() Avatar {
	nh.Lock()
	defer nh.Unlock()
	for _, av := range nh.open {
		return av
	}
	return Avatar{}
}
