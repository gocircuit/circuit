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

type Folk struct {
	kin *Kin
	topic string
	sync.Mutex
	ch chan FolkXID // Services pending to be opened
}

func (folk *Folk) Opened() []FolkXID {
	???
	neighbors := folk.kin.Neighbors()
	r := make([]FolkXID, len(neighbors))
	for i, v := range neighbors {
		r[i] = FolkXID(v)
	}
	return r
}

// Replenish blocks and returns the next downstream peer when one is chosen by the kin system.
func (folk *Folk) Replenish() (peer FolkXID) {
	return <-folk.ch
}

func (folk *Folk) supply(peer FolkXID) {
	if XID(peer).IsNil() {
		return
	}
	folk.Lock()
	defer folk.Unlock()
	folk.ch <- peer
}
