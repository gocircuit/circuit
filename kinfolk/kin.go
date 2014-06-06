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
		ach:   make(chan KinXID, 3*Spread),
		dch:   make(chan KinXID, 3*Spread),
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
	}() // XXX: Fewer than VertexExpansion entries triggers expand up
	ykin := YKin{KinXID{X: circuit.Dial(join, ServiceName)}}
	var w bytes.Buffer
	for _, peer := range ykin.Join() {
		peer = k.attach(peer)
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
		func (r interface{}) {
			k.rtr.Scrub(peer)
			k.expand()
		},
	))
	k.rtr.Add(peer)
	k.ach <- peer
	return peer
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
		kinXID := XKin{k}.Walk(Depth)
		if XID(kinXID).IsNil() || kinXID.ID == k.XID().ID { // Compare just IDs, in case we got pointers to ourselves from elsewhere
			continue // If peer is nil or self, ignore it
		}
		kinXID = k.remember(kinXID)
		log.Printf("expanding kinfolk system with a random kin %s", kinXID.ID.String())
		for _, folk := range k.users() {
			folk.supply(YKin{kinXID}.Attach(folk.Topic))
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

func (k *Kin) shrink() {
	for i := 0; i < k.rtr.Len() - ExpansionHigh; i++ {
		k.rtr.ScrubRandom()
	}
}

??

func (k *Kin) scrub(kinXID KinXID) {
	//debug.PrintStack()
	if k.rtr.Scrub(kinXID.X).IsNil() {
		return
	}
	log.Printf("Scrubbing kin %s", kinXID.ID.String())
	k.expand()
	k.dch <- kinXID
}

func (k *Kin) Attach(topic string, folkXID FolkXID) *Folk {
	// Fetch initial peer connections
	var opened = k.rtr.Opened()
	peers := make([]FolkXID, 0, len(opened))
	for _, xid := range opened {
		kinXID := KinXID(xid)
		if folkXID := (YKin{kinXID}).Attach(topic); folkXID.X != nil {
			// Scrub kinXID if folkXID.X ever panics
			folkXID.X = watch(folkXID.X, func(_ circuit.PermX, r interface{}) {
				k.scrub(kinXID)
				panic(r)
			})
			peers = append(peers, folkXID)
		}
	}

	k.Lock()
	defer k.Unlock()

	if _, present := k.topic[topic]; present {
		panic("dup attach")
	}
	folk := &Folk{
		Topic: topic,
		rtr:   NewRotor(),
		ch:    make(chan FolkXID, 6*(1+len(peers))),
	}
	// Inject initial peers
	for _, peer := range peers {
		folk.supply(peer)
	}

	// Register
	k.folk = append(k.folk, folk)
	k.topic[topic] = folkXID

	return folk
}
