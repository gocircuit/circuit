// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package circuit

import (
	"fmt"
	"math/rand"

	"github.com/hoijui/circuit/pkg/use/n"
)

// HandleID is a universal ID referring to a circuit value accessible across workers
type HandleID uint64

func (h HandleID) String() string {
	return fmt.Sprintf("H%016x", uint64(h))
}

// ChooseHandleID returns a random ID
func ChooseHandleID() HandleID {
	return HandleID(rand.Int63())
}

// X represents a cross-interface value.
type X interface {

	// Addr returns the address of the runtime, hosting the object underlying the cross-interface value.
	Addr() n.Addr

	// HandleID uniquely identifies the local reference to the receiver that was exported for this cross-reference
	HandleID() HandleID

	// Call invokes the method named proc of the actual object (possibly
	// living remotely) underlying the cross-interface. The invokation
	// arguments are take from in, and the returned values are placed in
	// the returned slice.
	//
	// Errors can only occur as a result of physical/external circumstances
	// that impede cross-worker communication. Such errors are returned in
	// the form of panics.
	Call(proc string, in ...interface{}) []interface{}

	// IsX is used internally.
	IsX()

	// String returns a human-readable representation of the cross-interface.
	String() string
}

// PermX represents a permanent cross-interface value.
type PermX interface {

	// A permanent cross-interface can be used as a non-permanent one.
	X

	// IsPerm is used internally.
	IsPermX()
}

// Func is a symbolic type that refers to circuit worker function types.
// These are types with a singleton public method.
type Func interface{}
