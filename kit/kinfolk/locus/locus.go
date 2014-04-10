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

	"github.com/gocircuit/circuit/kit/kinfolk"
	"github.com/gocircuit/circuit/kit/kinfolk/tube"
	nsrv "github.com/gocircuit/circuit/kit/fs/namespace/server"
	"github.com/gocircuit/circuit/kit/fs/client"
	"github.com/gocircuit/circuit/kit/fs/server"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/rh/xy"
	"github.com/gocircuit/circuit/use/circuit"
)

// Resources captures the local OS resources, shared with other workers
type Resources struct {
	Kin               *kinfolk.Kin            // The kinfolk social system of this circuit worker
	KinJoin, KinLeave <-chan kinfolk.KinXID   // Channels where peer joins and leaves are announced
	TubeTopic         string                  // Kinfolk tube topic used by the locus system
	FS                rh.Server               // Shared local file system
}

// Locus is a device that listens to the join/leave events reported by the kinfolk social
// system, and maintains an asynchronously-readable current list of known peers.
type Locus struct {
	Peer       *client.Peer      // Client peer enclosure for this circuit locus
	rsc        *Resources        // Local shared resources
	tube       *tube.Tube        // Kinfolk broadcasting system
	serverDir  *server.ServerDir // Local RH server FID (not for cross-sharing purposes)
}

// NewLocus creates a new locus device.
func NewLocus(rsc *Resources) *Locus {
	var err error
	locus := &Locus{
		rsc:  rsc,
		tube: tube.New(rsc.Kin, rsc.TubeTopic),
	}
	// Create an RH directory lens into the local file system, to give to the RH server namespace
	var fsFID rh.FID
	if rsc.FS != nil {
		fsSSN, err := rsc.FS.SignIn("", "/")
		if err != nil {
			panic(err)
		}
		if fsFID, err = fsSSN.Walk(nil); err != nil {
			panic(err)
		}
	}
	// Create server RH
	if fsFID != nil {
		locus.serverDir, err = server.NewDir(rsc.Kin.XID(), server.Resource{"fs", fsFID})
	} else {
		locus.serverDir, err = server.NewDir(rsc.Kin.XID())
	}
	if err != nil {
		panic(err)
	}
	//
	locus.Peer = &client.Peer{
		// It is crucial to use permanent cross-references, and not
		// "plain" ones within values stored inside the tube table. If
		// cross-references are used, they are managed by the cross-
		// garbage collection system and therefore connections to ALL
		// underlying workers are maintained superfluously.
		Kin:    rsc.Kin.XID(),
		Server: circuit.PermRef(xy.XServer{nsrv.NewServer(locus.serverDir)}),
	}
	//
	go loopJoin(rsc.KinJoin)
	go locus.loopLeave(rsc.KinLeave)
	go locus.loopExpire()
	//
	log.Println(locus.Peer.Key())
	return locus
}

func (locus *Locus) ServerDir() *server.ServerDir {
	return locus.serverDir
}

// GetPeers asynchronously returns the current known list of live peers.
func (locus *Locus) GetPeers() []*client.Peer {
	rr := locus.tube.BulkRead()
	s := make([]*client.Peer, len(rr))
	for i, r := range rr {
		s[i] = r.Value.(*client.Peer)
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
	peer := &client.Peer{Kin: kinXID}
	log.Println("denouncing", peer.Key())
	r := locus.tube.Lookup(peer.Key())
	if r == nil {
		return
	}
	locus.tube.Scrub(peer.Key(), r.Rev, r.Updated)
}
