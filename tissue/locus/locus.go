// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package locus

import (
	"log"
	"path"
	"time"

	"github.com/gocircuit/circuit/anchor"
	srv "github.com/gocircuit/circuit/element/server"
	"github.com/gocircuit/circuit/kit/pubsub"
	"github.com/gocircuit/circuit/tissue"
	"github.com/gocircuit/circuit/tissue/tube"
	"github.com/gocircuit/circuit/use/circuit"
)

// Locus is a device that listens to the join/leave events reported by the tissue social
// system, and maintains an asynchronously-readable current list of known peers.
type Locus struct {
	Peer *Peer      // Client peer enclosure for this circuit locus
	tube *tube.Tube // Kinfolk broadcasting system
}

// NewLocus creates a new locus device.
func NewLocus(kin *tissue.Kin, rip <-chan tissue.KinAvatar) XLocus {
	locus := &Locus{
		tube: tube.NewTube(kin, "locus"),
	}
	term, xterm := anchor.NewTerm(kin.Avatar().ID.String(), locus)
	term.Attach(anchor.Server, srv.New(kin))
	locus.Peer = &Peer{
		// It is crucial to use permanent cross-references, and not
		// "plain" ones within values stored inside the tube table. If
		// cross-references are used, they are managed by the cross-
		// garbage collection system and therefore connections to ALL
		// underlying workers are maintained superfluously.
		Kin:  kin.Avatar(),
		Term: xterm,
	}
	go locus.loopRIP(rip)
	go locus.loopAnnounceAndExpire()
	return XLocus{locus}
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

func (locus *Locus) Self() *Peer {
	return locus.Peer
}

// peerSubscription
type peerSubscription struct {
	pubsub.Consumer
}

func init() {
	circuit.RegisterValue(&peerSubscription{})
}

func (a *peerSubscription) X() circuit.X {
	return circuit.Ref(a)
}

func (a *peerSubscription) Consume() (interface{}, bool) {
	v, ok := a.Consumer.Consume()
	if !ok {
		return nil, false
	}
	return path.Join("/", v.(*tube.Record).Key), true
}

func (locus *Locus) NewArrivals() pubsub.Consumer {
	return &peerSubscription{locus.tube.NewArrivals()}
}

func (locus *Locus) NewDepartures() pubsub.Consumer {
	return &peerSubscription{locus.tube.NewDepartures()}
}

// loopAnnounceAndExpire writes a new version of this server's peer record to the tube view every 2 seconds,
// then it iterates through all peer records in the tube view, forgetting those older than 4 seconds.
func (locus *Locus) loopAnnounceAndExpire() {
	const GarbageDuration = time.Second * 4
	var rev tube.Rev
	for {
		rev++
		// log.Printf("(Re)announcing ourselves (%s,%d,%v)", locus.Peer.Key(), rev, locus.Peer)
		locus.tube.Write(locus.Peer.Key(), rev, locus.Peer)
		//
		time.Sleep(GarbageDuration / 2)
		deadline := time.Now().Add(-GarbageDuration)
		for _, r := range locus.tube.BulkRead() {
			locus.tube.Forget(r.Key, 0, deadline)
		}
	}
}

func (locus *Locus) loopRIP(rip <-chan tissue.KinAvatar) {
	for {
		kinAvatar, ok := <-rip
		if !ok {
			panic("u")
		}
		locus.denounce(kinAvatar)
	}
}

func (locus *Locus) denounce(kinAvatar tissue.KinAvatar) {
	peer := &Peer{Kin: kinAvatar}
	log.Println("Denouncing", peer.Key())
	r := locus.tube.Lookup(peer.Key())
	if r == nil {
		return
	}
	locus.tube.Scrub(peer.Key(), r.Rev, r.Updated)
}
