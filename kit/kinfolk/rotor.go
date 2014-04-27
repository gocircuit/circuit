// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package kinfolk

import (
	"math/rand"
	"sync"

	"github.com/gocircuit/circuit/use/circuit"
)

// Rotor maintains a set of XIDs (cross-interfaces, paired with unique
// receiver identifiers). When an XID is added to the rotor with Open, it is
// equipped with a watchdog logic that automatically removes the XID from the
// rotor, if any invocation to its Call method panics.
type Rotor struct {
	sync.Mutex
	open []XID
}

// NewRotor creates a new rotor.
func NewRotor() *Rotor {
	return &Rotor{
		open: make([]XID, 0, 6*Spread),
	}
}

func (rtr *Rotor) add(xid XID) {
	rtr.Lock()
	defer rtr.Unlock()
	rtr.open = append(rtr.open, xid)
}

// Scrub looks up an XID by the value of its cross-interface;
// If it finds a matching XID, it removes it from the rotor
// and returns the removed XID.
func (rtr *Rotor) Scrub(x circuit.X) XID {
	rtr.Lock()
	defer rtr.Unlock()
	for i, xid := range rtr.open {
		if xid.X == x {
			n := len(rtr.open)
			rtr.open[i] = rtr.open[n-1]
			rtr.open = rtr.open[:n-1]
			return xid
		}
	}
	return XID{}
}

// Opened returns a list of all open and healthy XIDs
func (rtr *Rotor) Opened() []XID {
	rtr.Lock()
	defer rtr.Unlock()
	open := make([]XID, len(rtr.open))
	for i, j := range rand.Perm(len(rtr.open)) {
		open[i] = rtr.open[j]
	}
	return open
}

// NOpened returns the number of XIDs in the rotor.
func (rtr *Rotor) NOpened() int {
	rtr.Lock()
	defer rtr.Unlock()
	return len(rtr.open)
}

// Choose returns a randomly chosen XID.
func (rtr *Rotor) Choose() XID {
	rtr.Lock()
	defer rtr.Unlock()
	n := len(rtr.open)
	if n <= 0 {
		return XID{}
	}
	return rtr.open[rand.Intn(n)]
}

// Open moves the cross-interface XID to the set of open XIDs.
func (rtr *Rotor) Open(xid XID) XID {
	scrubx := watch(xid.X, func(scrubx circuit.PermX, r interface{}) {
		rtr.Scrub(scrubx)
		panic(r)
	})
	xid = XID{X: scrubx, ID: xid.ID}
	rtr.add(xid)
	return xid
}
