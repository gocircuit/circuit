// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package tele implements the circuit/use/n networking module using Teleport Transport
package tele

import (
	"net"
	"os"

	"github.com/gocircuit/circuit/kit/tele"
	"github.com/gocircuit/circuit/kit/tele/blend"
	"github.com/gocircuit/circuit/use/n"
)

func init() {
	n.Bind(&System{})
}

// System is the high-level type that encloses a monolithic networking functionality
type System struct{}

// workerID is the ID for this transport endpoint.
// addr is the networking address to listen to.
func (s *System) NewTransport(workerID n.WorkerID, addr net.Addr, key []byte) n.Transport {
	var u *blend.Transport
	if len(key) == 0 {
		u = tele.NewStructOverTCP()
	} else {
		u = tele.NewStructOverTCPWithHMAC(key)
	}
	l := newListener(workerID, os.Getpid(), u.Listen(addr))
	return &Transport{
		WorkerID: workerID,
		Dialer:   newDialer(l.Addr(), u),
		Listener: l,
	}
}

func (s *System) ParseNetAddr(a string) (net.Addr, error) {
	return ParseNetAddr(a)
}

func (s *System) ParseAddr(a string) (n.Addr, error) {
	return ParseAddr(a)
}

// Transport cumulatively represents the ability to listen for connections and dial into remote endpoints.
type Transport struct {
	n.WorkerID
	*Dialer
	*Listener
}
