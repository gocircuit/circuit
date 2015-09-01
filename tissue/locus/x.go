// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package locus

import (
	"github.com/gocircuit/circuit/use/circuit"
)

func init() {
	circuit.RegisterValue(XLocus{})
}

type XLocus struct {
	l *Locus
}

func (x XLocus) GetPeers() []*Peer {
	return x.l.GetPeers()
}

func (x XLocus) Self() interface{} {
	return x.l.Self()
}

// YLocus…
type YLocus struct {
	X circuit.PermX
}

func (y YLocus) GetPeers() map[string]*Peer {
	r := make(map[string]*Peer)
	for _, p := range y.X.Call("GetPeers")[0].([]*Peer) {
		r[p.Key()] = p
	}
	return r
}

func (y YLocus) Self() *Peer {
	return y.X.Call("Self")[0].(*Peer)
}
