// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tissue

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
	// for the purposes of dynamically maintaining node presence using the tissue
	// collaborative protocol.
	ExpansionLow  = 7
	ExpansionHigh = 11
	Spread        = 5

	// Depth is the number of random walk steps taken when sampling for a random circuit worker.
	Depth = 3 * 2 // Lazy random walk with stay-put probability one half
)

// Avatar is a pair of a permanent cross-interface and an ID, identifying its underlying receiver uniquely.
type Avatar struct {
	X  circuit.PermX
	ID lang.ReceiverID
}

func (av Avatar) Equals(w Avatar) bool {
	if av.ID == 0 {
		panic(0)
	}
	return av.ID == w.ID
}

// String returns a textual form of the Avatar.
func (av Avatar) String() string {
	return fmt.Sprintf("%s ==> %s", av.ID, av.X.Addr())
}

// IsNil returns true if the Avatar is not set.
func (av Avatar) IsNil() bool {
	return av.X == nil
}

// ForwardPanic replaces the cross-interface x with one that
// captures panics during method calls and passes them to the fwd func in a separate goroutine,
// while also propagating the original panic through the stack.
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

// ForwardAvatarPanic is like ForwardPanic but for Avatar objects.
func ForwardAvatarPanic(av Avatar, fwd func(recov interface{})) Avatar {
	av.X = ForwardPanic(av.X, fwd)
	return av
}
