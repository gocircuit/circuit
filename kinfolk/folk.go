// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package kinfolk

import (
	"sync"
)

//
type Folk struct {
	Topic string
	rtr   *Rotor
	sync.Mutex
	ch chan FolkXID // Services pending to be opened
}

//
func (folk *Folk) Opened() []FolkXID {
	o := folk.rtr.Opened()
	r := make([]FolkXID, len(o))
	for i, v := range o {
		r[i] = FolkXID(v)
	}
	return r
}

// Replenish blocks and returns the next downstream peer when one is chosen by the kin system.
func (folk *Folk) Replenish() (peer FolkXID) {
	x := <-folk.ch
	return FolkXID(folk.rtr.Open(XID(x)))
}

func (folk *Folk) supply(peer FolkXID) {
	if XID(peer).IsNil() {
		return
	}
	folk.Lock()
	defer folk.Unlock()
	folk.ch <- peer
}
