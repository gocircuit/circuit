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

// Init
func init() {
	circuit.RegisterValue(&Kin{}) // So that we can computer receiver ID
	circuit.RegisterValue(ExoKin{})
}

const (
	Spread = 5
	Depth  = 5
)

// XID is a pair of a cross-interface and an ID, identifying its receiver uniquely
type XID struct {
	X  circuit.PermX
	ID lang.ReceiverID
}

func (xid XID) String() string {
	return fmt.Sprintf("%s ==> %s", xid.ID, xid.X.Addr())
}

func (xid XID) IsNil() bool {
	return xid.X == nil
}

// watch returns an equivalent cross-reference to x, which will execute the function r
// each time an invokation to x.Call results in panic.
func watch(x circuit.PermX, recov func(wx circuit.PermX, r interface{})) circuit.PermX {
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
