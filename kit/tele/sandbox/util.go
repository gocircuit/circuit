// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package sandbox provides a simulated carrier Transport for testing purposes.
package sandbox

import (
	"math/rand"
	"net"
	"time"

	"github.com/gocircuit/circuit/kit/tele/trace"
)

func NewReliableTransport(f trace.Frame) *Transport {
	return NewTransport(f, NewPipe)
}

func NewUnreliableTransport(f trace.Frame, nok, ndrop int, expa, expb time.Duration) *Transport {
	return NewTransport(f, func(f0, f1 trace.Frame, a0, a1 net.Addr) (net.Conn, net.Conn) {
		f.Printf("TRANSPORT PROFILE NOK=%d, NDROP=%d", nok, ndrop)
		return NewSievePipe(f0, f1, a0, a1, nok, ndrop, expa, expb)
	})
}

func NewRandomUnreliableTransport(f trace.Frame, nok, ndrop int, expa, expb time.Duration) *Transport {
	return NewTransport(f, func(f0, f1 trace.Frame, a0, a1 net.Addr) (net.Conn, net.Conn) {
		nok, ndrop := rand.Intn(nok+1), rand.Intn(ndrop+1)
		nok = max(nok, 1)
		f.Printf("TRANSPORT PROFILE NOK=%d, NDROP=%d", nok, ndrop)
		return NewSievePipe(f0, f1, a0, a1, nok, ndrop, expa, expb)
	})
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
