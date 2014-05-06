// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package locus

import (
	"log"
	"time"

	"github.com/gocircuit/circuit/kit/anchor"
	"github.com/gocircuit/circuit/kit/kinfolk"
	"github.com/gocircuit/circuit/kit/kinfolk/tube"
)

// Locus is a device that listens to the join/leave events reported by the kinfolk social
// system, and maintains an asynchronously-readable current list of known peers.
type Locus struct {
	Peer  *Peer      // Client peer enclosure for this circuit locus
	tube  *tube.Tube        // Kinfolk broadcasting system
}

// NewLocus creates a new locus device.
func NewLocus(kin *kinfolk.Kin, kinJoin, kinLeave <-chan kinfolk.KinXID) *Locus {
	locus := &Locus{
		tube: tube.New(kin, "locus"),
		Peer: &Peer{
			// It is crucial to use permanent cross-references, and not
			// "plain" ones within values stored inside the tube table. If
			// cross-references are used, they are managed by the cross-
			// garbage collection system and therefore connections to ALL
			// underlying workers are maintained superfluously.
			Kin:    kin.XID(),
			Term: anchor.NewTerm(kin.XID().ID.String()),
		},
	}
	go loopJoin(kinJoin)
	go locus.loopLeave(kinLeave)
	go locus.loopExpire()
	log.Println(locus.Peer.Key())
	return locus
}

// GetPeers asynchronously returns the current known list of live peers.
func (locus *Locus) GetPeers() []*Peer {
	rr := locus.tube.BulkRead()
	s := make([]*Peer, len(rr))
	for i, r := range rr {
		s[i] = r.Value.(*Peer)
	}
	return s
}

func (locus *Locus) loopExpire() {
	const GarbageDuration = time.Second * 3
	var rev tube.Rev
	for {
		rev++
		//log.Printf("WRITING (%s,%d,%v)", locus.info.Key(), rev, locus.info)
		locus.tube.Write(locus.Peer.Key(), rev, locus.Peer)
		//
		time.Sleep(GarbageDuration / 2)
		deadline := time.Now().Add(-GarbageDuration)
		for _, r := range locus.tube.BulkRead() {
			locus.tube.Forget(r.Key, 0, deadline)
		}
	}
}

func loopJoin(kinjoin <-chan kinfolk.KinXID) {
	// Discard join events
	for {
		kinXID, ok := <-kinjoin
		if !ok {
			panic("u")
		}
		log.Println("kin joined", kinXID.String())
	}
}

func (locus *Locus) loopLeave(kinleave <-chan kinfolk.KinXID) {
	for {
		kinXID, ok := <-kinleave
		if !ok {
			panic("u")
		}
		locus.denounce(kinXID)
	}
}

func (locus *Locus) denounce(kinXID kinfolk.KinXID) {
	peer := &Peer{Kin: kinXID}
	log.Println("denouncing", peer.Key())
	r := locus.tube.Lookup(peer.Key())
	if r == nil {
		return
	}
	locus.tube.Scrub(peer.Key(), r.Rev, r.Updated)
}
