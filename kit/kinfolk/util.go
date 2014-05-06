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
	// Spread is the number of peers that each circuit worker is continuously connected to  
	// for the purposes of dynamically maintaining node presence using the kinfolk
	// collaborative protocol.
	Spread = 5

	// Depth is the number of random walk steps taken when sampling for a random circuit worker.
	Depth  = 5
)

// XID is a pair of a permanent cross-interface and an ID, identifying its underlying receiver uniquely.
type XID struct {
	X  circuit.PermX
	ID lang.ReceiverID
}

// String returns a textual form of the XID.
func (xid XID) String() string {
	return fmt.Sprintf("%s ==> %s", xid.ID, xid.X.Addr())
}

// IsNil returns true if the XID is not set.
func (xid XID) IsNil() bool {
	return xid.X == nil
}

// watch returns an equivalent cross-reference to x, which will execute the function r
// each time an invokation to x.Call results in panic.
type recoverFunc func(wx circuit.PermX, r interface{})

func watch(x circuit.PermX, recov recoverFunc) circuit.PermX {
	return &watchx{PermX: x, r: recov}
}

type watchx struct {
	circuit.PermX
	r func(circuit.PermX, interface{})
}

func (wx *watchx) Call(proc string, in ...interface{}) []interface{} {
	defer func() {
		if r := recover(); r != nil {
			wx.r(wx, r)
		}
	}()
	return wx.PermX.Call(proc, in...)
}
