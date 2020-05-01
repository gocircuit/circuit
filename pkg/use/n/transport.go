// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package n

import (
	"net"
)

const Scheme = "circuit"

// Addr represents the identity of a unique remote worker.
// The implementing type must be registered with package encoding/gob.
type Addr interface {
	// NetAddr returns the underlying networking address of this endpoint.
	NetAddr() net.Addr

	// String returns an equivalent textual representation of the address.
	String() string

	// Returns a unique textual representation of the address, suitable to be used as a file name.
	FileName() string

	// WorkerID returns the worker ID of the underlying worker.
	WorkerID() WorkerID
}

// Conn is a connection to a remote endpoint.
type Conn interface {

	// The language runtime does not itself utilize timeouts on read/write
	// operations. Instead, it requires that calls to Read and Write be blocking
	// until success or irrecoverable failure is reached.
	//
	// The implementation of Conn must monitor the liveness of the remote
	// endpoint using an out-of-band (non-visible to the runtime) method. If
	// the endpoint is considered dead, all pending Read and Write request must
	// return with non-nil error.
	//
	// A non-nil error returned on any invokation of Read and Write signals to
	// the runtime that not just the connection, but the entire runtime
	// (identified by its address) behind the connection is dead.
	//
	// Such an event triggers various language runtime actions such as, for
	// example, releasing all values exported to that runtime. Therefore, a
	// typical Conn implementation might choose to attempt various physical
	// connectivity recovery methods, before it reports an error on any pending
	// connection. Such implentation strategies are facilitated by the fact
	// that the runtime has no semantic limits on the length of blocking waits.
	// In fact, the runtime has no notion of time altogether.

	// Read/Write operations must panic on any encoding/decoding errors.
	// Whereas they must return an error for any exernal (network) unexpected
	// conditions.  Encoding errors indicate compile-time errors (that will be
	// caught automatically once the system has its own compiler) but might be
	// missed by the bare Go compiler.
	//
	// Read/Write must be re-entrant.

	// Read reads the next value from the connection.
	Read() (interface{}, error)

	// Write writes the given value to the connection.
	Write(interface{}) error

	// Close closes the connection.
	Close() error

	// ...
	Abort(reason error)

	// Addr returns the address of the remote endpoint.
	Addr() Addr

	// TODO: ReverseDial is potentially a better abstraction for the heartbeat connections created when passing cross-interfaces
	// ReverseDial() (Conn, error)
}

// Listener is a device for accepting incoming connections.
type Listener interface {

	// Accept returns the next incoming connection.
	Accept() Conn

	// Close closes the listening device.
	Close()

	// Addr returns the address of this endpoint.
	Addr() Addr
}

// Dialer is a device for initating connections to addressed remote endpoints.
type Dialer interface {

	// Dial connects to the endpoint specified by addr and returns a respective connection object.
	Dial(addr Addr) (Conn, error)
}

// Transport cumulatively represents the ability to listen for connections and dial into remote endpoints.
type Transport interface {
	Dialer
	Listener
}

// System creates a new transport framework for the given local address
type System interface {
	NewTransport(workerID WorkerID, addr net.Addr, key []byte) Transport
	ParseNetAddr(s string) (net.Addr, error)
	ParseAddr(s string) (Addr, error)
}
