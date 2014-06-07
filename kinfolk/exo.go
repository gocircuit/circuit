// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package kinfolk

import (
	"time"

	"github.com/gocircuit/circuit/use/n"
	"github.com/gocircuit/circuit/kit/lang"
)

// FolkXID is an XID underlied by a user receiver for a service shared over the kinfolk system.
type FolkXID XID

func (xid FolkXID) XID() XID {
	return XID(xid)
}

func (xid FolkXID) String() string {
	return "FolkXID:" + XID(xid).String()
}

// KinXID is an XID specifically for the XKin receiver.
type KinXID XID

func (xid KinXID) String() string {
	return "KinXID:" + XID(xid).String()
}

// XKin is the cross-worker interface of the kinfolk system at this circuit.
type XKin struct {
	k *Kin
}

// Attach returns a cross-reference to a folk service at this worker.
func (x XKin) Attach(topic string) FolkXID {
	x.k.Lock()
	defer x.k.Unlock()
	return x.k.topic[topic]
}

// Join returns an initial set of peers that the joining kin should use as initial entry into the system.
func (x XKin) Join(joinee n.Addr) []KinXID {
	m := make(map[lang.ReceiverID]KinXID)
	for i := 0; i < Spread; i++ {
		peerXID := x.Walk(Depth)
		if XID(peerXID).IsNil() {
			continue
		}
		if _, ok := m[peerXID.ID]; ok {
			// Duplicate
			continue
		}
		m[peerXID.ID] = peerXID
	}
	r := make([]KinXID, 0, len(m))
	for _, peerXID := range m {
		r = append(r, peerXID)
	}
	if joinee != nil {
		go func() {
			time.Sleep(time.Second/2)
			x.k.ReJoin(joinee) // Reciprocate by joining into their network
		}()
	}
	return r
}

// Walk performs a random walk through the expander-graph network of circuit workers
// of length t steps and returns the kinfolk XID of the terminal node.
func (x XKin) Walk(t int) KinXID {
	if t <= 0 {
		return x.k.XID()
	}
	hop := KinXID(x.k.neighborhood.Choose())
	if hop.X == nil {
		return x.k.XID()
	}
	defer func() {
		recover()
	}()
	return YKin{hop}.Walk(t - 1)
}

// YKin
type YKin struct {
	xid KinXID
}

func (y YKin) Join(joinee n.Addr) []KinXID {
	// Do not recover
	return y.xid.X.Call("Join", joinee)[0].([]KinXID)
}

func (y YKin) Walk(t int) KinXID {
	// Do not recover; XKin.Walk relies on panics
	return y.xid.X.Call("Walk", t)[0].(KinXID)
}

func (y YKin) Attach(topic string) FolkXID {
	defer func() {
		recover()
	}()
	return y.xid.X.Call("Attach", topic)[0].(FolkXID)
}
