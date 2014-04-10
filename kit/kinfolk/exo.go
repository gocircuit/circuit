// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package kinfolk ...
package kinfolk

import (
	"github.com/gocircuit/circuit/kit/lang"
	"github.com/gocircuit/circuit/use/circuit"
)

// XKin is a permanent cross-interface to ExoKin
type XKin circuit.PermX

// FolkXID is the XID of a …
type FolkXID XID

func (xid FolkXID) XID() XID {
	return XID(xid)
}

func (xid FolkXID) String() string {
	return "FolkXID." + XID(xid).String()
}

// KinXID is an XID specifically for the ExoKin receiver.
type KinXID XID

func (xid KinXID) String() string {
	return "KinXID." + XID(xid).String()
}

// ExoKin is the cross-worker interface of the kinfolk system at this circuit.
type ExoKin struct {
	k *Kin
}

// Attach returns a cross-reference to an exo folk object
func (x ExoKin) Attach(topic string) FolkXID {
	x.k.Lock()
	defer x.k.Unlock()
	return x.k.topic[topic]
}

// Join returns an initial set of peers that the joining kin should use as initial entry into the kin system.
func (x ExoKin) Join() []KinXID {
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
	return r
}

func (x ExoKin) Walk(t int) KinXID {
	if t <= 0 {
		return x.k.XID()
	}
	hop := KinXID(x.k.rtr.Choose())
	if hop.X == nil {
		return x.k.XID()
	}
	defer func() {
		if r := recover(); r != nil {
			x.k.scrub(hop)
		}
	}()
	return YKin{hop}.Walk(t - 1)
}

// YKin
type YKin struct {
	xid KinXID
}

func (y YKin) Join() []KinXID {
	// Do not recover
	return y.xid.X.Call("Join")[0].([]KinXID)
}

func (y YKin) Walk(t int) KinXID {
	// Do not recover; ExoKin.Walk relies on panics
	return y.xid.X.Call("Walk", t)[0].(KinXID)
}

func (y YKin) Attach(topic string) FolkXID {
	defer func() {
		recover()
	}()
	return y.xid.X.Call("Attach", topic)[0].(FolkXID)
}
