// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package kinfolk

import (
	"sync"

	"github.com/gocircuit/circuit/kit/lang"
)

// Neighborhood is a set of perm cross-interfaces.
type Neighborhood struct {
	sync.Mutex
	open map[interface{}]XID
}

// NewNeighborhood creates a new rotor.
func NewNeighborhood() *Neighborhood {
	return &Neighborhood{
		open: make(map[interface{}]XID),
	}
}

func (nh *Neighborhood) Add(xid XID) {
	nh.Lock()
	defer nh.Unlock()
	nh.open[xid.ID] = xid
}

func (nh *Neighborhood) Scrub(key lang.ReceiverID) (XID, bool) {
	nh.Lock()
	defer nh.Unlock()
	xid, ok := nh.open[key]
	delete(nh.open, key)
	return xid, ok
}

func (nh *Neighborhood) ScrubRandom() (XID, bool) {
	nh.Lock()
	defer nh.Unlock()
	for key, xid := range nh.open {
		delete(nh.open, key)
		return xid, true
	}
	return XID{}, false
}

// View returns a list of all XIDs in the rotor.
func (nh *Neighborhood) View() []XID {
	nh.Lock()
	defer nh.Unlock()
	open := make([]XID, 0, len(nh.open))
	for _, xid := range nh.open {
		open = append(open, xid)
	}
	return open
}

// Len returns the number of XIDs in the rotor.
func (nh *Neighborhood) Len() int {
	nh.Lock()
	defer nh.Unlock()
	return len(nh.open)
}

// Choose returns a randomly chosen XID.
func (nh *Neighborhood) Choose() XID {
	nh.Lock()
	defer nh.Unlock()
	for _, xid := range nh.open {
		return xid
	}
	return XID{}
}
