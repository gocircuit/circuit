// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"runtime"

	"github.com/gocircuit/circuit/kit/element/proc"
	"github.com/gocircuit/circuit/kit/element/valve"
	"github.com/gocircuit/circuit/use/circuit"
)

func init() {
	circuit.Register(XChamber{})
}

type XChamber struct {
	a Anchor
}

// Walk returns a cross-interface to a XChamber.
func (xmbr XChamber) Walk(walk []string) circuit.X {
	return circuit.Ref(XChamber{xmbr.a.Walk(walk)})
}

func (xmbr XChamber) View() (interface{}, map[string]circuit.X) {
	v, w := xmbr.a.View()
	u := make(map[string]circuit.X)
	for p, q := range w {
		u[p] = circuit.Ref(q)
	}
	return v, u
}

func (xmbr XChamber) MakeValve(n int) circuit.X {
	v := valve.MakeValve(n)
	if !xmbr.a.Set(v) {
		return nil
	}
	return circuit.Ref(valve.XValve{v})
}

func (xmbr XChamber) MakeProc(cmd *proc.Cmd) circuit.X { // XProc
	p := valve.MakeProc(cmd)
	if !xmbr.a.Set(p) {
		return nil
	}
	return circuit.Ref(valve.XValve{v})
}

func (xmbr XChamber) Scrub() {
	?
}

// YChamber …
type YChamber struct {
	x circuit.X
}

func (ymbr XChamber) Walk(walk []string) YChamber {
	?
}

func (ymbr XChamber) View() (interface{}, map[string]YChamber) { // interface is YValve or YProc
	?
}

func (ymbr XChamber) MakeValve(n int) YValve {
	?
}

func (ymbr XChamber) MakeProc(cmd *proc.Cmd) YProc {
	?
}

func (ymbr XChamber) Scrub() {
	?
}
