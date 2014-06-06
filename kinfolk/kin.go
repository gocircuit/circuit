// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package kinfolk is an efficient “social” protocol for maintaining mutual awareness and sharing resources amongs circuit workers.
package kinfolk

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	//"runtime/debug"

	"github.com/gocircuit/circuit/kit/lang"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
)

// Kin is the kinfolk system server.
type Kin struct {
	kinxid KinXID // Permanent circuit-wide unique ID for this kin
	rtr    *Rotor
	///
	ach    chan KinXID // announcements of newly discovered kins
	dch    chan KinXID // denouncements of newly discovered deceased kins
	///
	sync.Mutex
	topic  map[string]FolkXID // topic -> yfolk
	folk   []*Folk
}

const ServiceName = "kin"

func NewKin() (k *Kin, xkin XKin, join, leave <-chan KinXID) {
	k = &Kin{
		rtr:   NewRotor(),
		ach:   make(chan KinXID, ExpansionHigh),
		dch:   make(chan KinXID, ExpansionHigh),
		topic: make(map[string]FolkXID),
	}
	// Create a KinXID for this system.
	k.kinxid = KinXID(XID{
		X:  circuit.PermRef(XKin{k}),
		ID: lang.ComputeReceiverID(k),
	})
	return k, XKin{k}, k.ach, k.dch
}

// ReJoin contacts the peering kin service join and joins its circuit network.
func (k *Kin) ReJoin(join n.Addr) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic joining: %v", r)
		}
	}()
	ykin := YKin{KinXID{X: circuit.Dial(join, ServiceName)}}
	var w bytes.Buffer
	for _, peer := range ykin.Join() {
		peer = k.remember(peer)
		w.WriteString(peer.X.Addr().WorkerID().String())
		w.WriteByte(' ')
	}
	if w.Len() > 0 {
		log.Println("Adding peering server(s):", w.String())
	}
	return nil
}

func (k *Kin) remember(peer KinXID) KinXID {
	peer = KinXID(ForwardXIDPanic(XID(peer),
		func (interface{}) {
			k.forget(peer)
		},
	))
	k.rtr.Add(peer)
	k.ach <- peer
	k.shrink()
	return peer
}

// shrink shrinks the neighborhood down to size ExpansionHigh.
func (k *Kin) shrink() {
	for i := 0; i < k.rtr.Len() - ExpansionHigh; i++ {
		xid, ok := k.rtr.ScrubRandom()
		if !ok {
			return
		}
		log.Printf("Shrunk kin %s", xid.ID.String())
		k.dch <- KinXID(xid)
	}
}

func (k *Kin) forget(kinXID KinXID) {
	if !k.rtr.Scrub(kinXID) {
		return
	}
	log.Printf("Forgetting kin %s", kinXID.ID.String())
	k.expand()
	k.dch <- kinXID
}

func (k *Kin) XID() KinXID {
	return k.kinxid
}

func (k *Kin) expand() {
	if k.rtr.Len() < ExpansionLow {
		return
	}
	for i := 0; i < ExpansionHigh - k.rtr.Len(); i++ {
		// Choose a random peer, using a random walk
		w := XKin{k}.Walk(Depth)
		if XID(w).IsNil() || w.ID == k.XID().ID { // Compare just IDs, in case we got pointers to ourselves from elsewhere
			continue // If peer is nil or self, ignore it
		}
		w = k.remember(w)
		log.Printf("expanding kinfolk system with a random kin %s", w.ID.String())
		for _, folk := range k.users() {
			folk.supply(YKin{w}.Attach(folk.Topic))
		}
	}
}

func (k *Kin) users() []*Folk {
	k.Lock()
	defer k.Unlock()
	folk := make([]*Folk, len(k.folk))
	copy(folk, k.folk)
	return folk
}

func (k *Kin) Attach(topic string, folkXID FolkXID) *Folk {
	var neighbors = k.rtr.View()
	peers := make([]FolkXID, 0, len(neighbors))
	for _, xid := range neighbors {
		kinXID := KinXID(xid)
		if folkXID := (YKin{kinXID}).Attach(topic); folkXID.X != nil {
			peers = append(peers, folkXID)
		}
	}
	k.Lock()
	defer k.Unlock()
	if _, present := k.topic[topic]; present {
		panic("dup attach")
	}
	folk := &Folk{
		topic: topic,
		kin: k,
		ch:    make(chan FolkXID, 2*ExpansionHigh),
	}
	for _, peer := range peers {
		folk.supply(peer)
	}
	k.folk = append(k.folk, folk)
	k.topic[topic] = folkXID
	return folk
}
