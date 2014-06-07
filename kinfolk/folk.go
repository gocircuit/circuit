// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package kinfolk

import (
	//"sync"
)

type Folk struct {
	kin *Kin
	topic string
	neighborhood *Neighborhood
	ch chan FolkXID // Services pending to be opened
}

func (folk *Folk) Opened() []FolkXID {
	neighbors := folk.neighborhood.View()
	r := make([]FolkXID, len(neighbors))
	for i, v := range neighbors {
		r[i] = FolkXID(v)
	}
	return r
}

// Replenish blocks and returns the next downstream peer added to the neighborhod set by the kin.
func (folk *Folk) Replenish() (peer FolkXID) {
	peer = <-folk.ch
	folk.neighborhood.Add(XID(peer))
	return peer
}

func (folk *Folk) addPeer(peer FolkXID) {
	if XID(peer).IsNil() {
		return
	}
	folk.ch <- peer
}

func (folk *Folk) removePeer(peer FolkXID) {
	folk.neighborhood.Scrub(XID(peer))
}
