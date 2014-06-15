// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tissue

import (
	"bytes"
	//"log"
	"math/rand"
)

// FolkAvatar is an Avatar underlied by a user receiver for a service shared over the tissue system.
type FolkAvatar Avatar

func (av FolkAvatar) Avatar() Avatar {
	return Avatar(av)
}

func (av FolkAvatar) String() string {
	return "FolkAvatar:" + Avatar(av).String()
}

// KinAvatar is an Avatar specifically for the XKin receiver.
type KinAvatar Avatar

func (av KinAvatar) String() string {
	return "KinAvatar:" + Avatar(av).String()
}

// XKin is the cross-worker interface of the tissue system at this circuit.
type XKin struct {
	k *Kin
}

// Attach returns a cross-reference to a folk service at this worker.
func (x XKin) Attach(topic string) FolkAvatar {
	x.k.Lock()
	defer x.k.Unlock()
	return x.k.topic[topic]
}

// Join …
func (x XKin) Join(boundary []KinAvatar, spread int) []KinAvatar {
	offer := x.k.chooseBoundary(spread) // compute boundary before merge happens
	var w bytes.Buffer
	for _, q := range boundary {
		q = x.k.remember(q)
		w.WriteString(q.X.Addr().String())
		w.WriteByte(' ')
	}
	// if len(boundary) > 0 {
	// 	log.Println("Remembering merging server(s):", w.String())
	// }
	return offer
}

// Walk performs a random walk through the expander-graph network of circuit workers
// of length t steps and returns the tissue Avatar of the terminal node.
func (x XKin) Walk(t int) KinAvatar {
	if t <= 0 {
		return x.k.Avatar()
	}
	if rand.Intn(2) < 1 { // Lazy random walk
		return x.Walk(t-1)
	}
	hop := KinAvatar(x.k.neighborhood.Choose())
	if hop.X == nil {
		return x.k.Avatar()
	}
	defer func() {
		recover()
	}()
	return YKin{hop}.Walk(t - 1)
}

// YKin
type YKin struct {
	av KinAvatar
}

func (y YKin) Join(boundary []KinAvatar, n int) []KinAvatar {
	// Do not recover
	return y.av.X.Call("Join", boundary, n)[0].([]KinAvatar)
}

func (y YKin) Walk(t int) KinAvatar {
	// Do not recover; XKin.Walk relies on panics
	return y.av.X.Call("Walk", t)[0].(KinAvatar)
}

func (y YKin) Attach(topic string) FolkAvatar {
	defer func() {
		recover()
	}()
	return y.av.X.Call("Attach", topic)[0].(FolkAvatar)
}
