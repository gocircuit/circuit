// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tissue

import (
	//"sync"

	"github.com/hoijui/circuit/kit/lang"
)

type Folk struct {
	kin *Kin
	topic string
	neighborhood *Neighborhood
	ch chan FolkAvatar // Services pending to be opened
}

func (folk *Folk) Opened() []FolkAvatar {
	neighbors := folk.neighborhood.View()
	r := make([]FolkAvatar, len(neighbors))
	for i, v := range neighbors {
		r[i] = FolkAvatar(v)
	}
	return r
}

// Replenish blocks and returns the next downstream peer added to the neighborhod set by the kin.
func (folk *Folk) Replenish() (peer FolkAvatar) {
	peer = <-folk.ch
	folk.neighborhood.Add(Avatar(peer))
	return peer
}

func (folk *Folk) addPeer(peer FolkAvatar) {
	if Avatar(peer).IsNil() {
		return
	}
	folk.ch <- peer
}

func (folk *Folk) removePeer(key lang.ReceiverID) {
	folk.neighborhood.Scrub(key)
}
