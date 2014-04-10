// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chain

import (
	"errors"
	"net"
)

// Carrier represents a transport underlying the chain transport layer.
// Carrier is a reliable, connection-oriented transport akin to TCP.
type Carrier interface {

	// Listen returns a new listener that listens on the given opaque address.
	Listen(net.Addr) (net.Listener, error)

	// Dial tries to establish a connection with the addressed remote endpoint.
	// Dial must distinguish between temporary and permanent obstructions to dialing the destination.
	// An ErrRIP error should be returned if addressed entity is permanently (from a logical standpoint) dead.
	// All other errors are considered temporary obstructions.
	Dial(net.Addr) (net.Conn, error)
}

var ErrRIP = errors.New("remote permanently gone")
