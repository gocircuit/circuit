// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package kinfolk

import (
	"fmt"

	"github.com/gocircuit/circuit/kit/lang"
	"github.com/gocircuit/circuit/use/circuit"
)

func init() {
	circuit.RegisterValue(&Kin{}) // So that we can compute receiver ID
	circuit.RegisterValue(XKin{})
}

const (
	// Expansion is the number of peers that each circuit worker is continuously connected to  
	// for the purposes of dynamically maintaining node presence using the kinfolk
	// collaborative protocol.
	ExpansionLow = 7
	ExpansionHigh = 11
	Spread = 5

	// Depth is the number of random walk steps taken when sampling for a random circuit worker.
	Depth  = 5
)

// XID is a pair of a permanent cross-interface and an ID, identifying its underlying receiver uniquely.
type XID struct {
	X  circuit.PermX
	ID lang.ReceiverID
}

func (xid XID) Equals(w XID) bool {
	if xid.ID == 0 {
		panic(0)
	}
	return xid.ID == w.ID
}

// String returns a textual form of the XID.
func (xid XID) String() string {
	return fmt.Sprintf("%s ==> %s", xid.ID, xid.X.Addr())
}

// IsNil returns true if the XID is not set.
func (xid XID) IsNil() bool {
	return xid.X == nil
}

// ForwardPanic…
func ForwardPanic(x circuit.PermX, fwd func(recov interface{})) circuit.PermX {
	return &forwardPanic{PermX: x, fwd: fwd}
}

type forwardPanic struct {
	circuit.PermX
	fwd func(interface{})
}

func (fp *forwardPanic) Call(proc string, in ...interface{}) []interface{} {
	defer func() {
		if r := recover(); r != nil {
			go fp.fwd(r)
			panic(r) // panics still have to go to the user
		}
	}()
	return fp.PermX.Call(proc, in...)
}

// ForwardXIDPanic…
fun ForwardXIDPanic(xid XID, fwd func(recov interface{})) XID {
	xid.X = ForwardPanic(xid.X, recov)
	return xid
}
