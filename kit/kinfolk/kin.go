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
	"log"
	"sync"
	//"runtime/debug"

	"github.com/gocircuit/circuit/kit/lang"
	"github.com/gocircuit/circuit/use/circuit"
)

// Kin is the kinfolk system logic, visible inside the circuit process
type Kin struct {
	kinxid KinXID // Permanent circuit-wide unique ID for this kin
	rtr    *Rotor
	ach    chan KinXID // announcements of newly discovered kins
	dch    chan KinXID // denouncements of newly discovered deceased kins
	sync.Mutex
	topic  map[string]FolkXID // topic -> yfolk
	folk   []*Folk
}

// NewKin creates a new kinfolk system server which, optionally,
// joins the kinfolk network that join is a member of; join is a
// permanent cross-interface to a peering kinfolk system.
//
// Additions and removals of new members to the circuit system
// will be announced over add and rmv. These two channels must 
// be consumed by the user.
func NewKin(join circuit.PermX) (k *Kin, xkin XKin, add, rmv <-chan KinXID) {
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
	if join == nil {
		return k, XKin{k}, k.ach, k.dch
	}
	// Join peer
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("joining (%s)", r)
		}
	}()
	var w bytes.Buffer
	for _, peer := range (YKin{ KinXID{ X: join }}).Join() {
		k.open(peer)
		w.WriteString(peer.X.Addr().WorkerID().String())
		w.WriteByte(' ')
	}
	log.Println("Initial peering servers:", w.String())
	return k, XKin{k}, k.ach, k.dch
}

func (k *Kin) open(peer KinXID) KinXID {
	peer = KinXID(k.rtr.Open(XID(peer)))
	k.ach <- peer
	return peer
}

func (k *Kin) XID() KinXID {
	return k.kinxid
}

func (k *Kin) replenish() {
	if k.rtr.NOpened() >= Spread {
		return
	}
	// Walk to a new XKin peer
	kinXID := XKin{k}.Walk(Depth)
	if XID(kinXID).IsNil() || kinXID.ID == k.XID().ID { // Compare just IDs, in case we got pointers to ourselves from elsewhere
		// If peer is nil or self, ignore it
		return
	}
	// Open the peer
	kinXID = k.open(kinXID)

	// Send new peer to all folk
	for _, folk := range k.snapfolk() {
		folk.supply(YKin{kinXID}.Attach(folk.Topic))
	}
}

func (k *Kin) snapfolk() []*Folk {
	k.Lock()
	defer k.Unlock()
	folk := make([]*Folk, len(k.folk))
	copy(folk, k.folk)
	return folk
}

func (k *Kin) scrub(kinXID KinXID) {
	//debug.PrintStack()
	if k.rtr.Scrub(kinXID.X).IsNil() {
		return
	}
	k.replenish()
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
