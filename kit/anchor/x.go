// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"errors"

	"github.com/gocircuit/circuit/kit/element/proc"
	"github.com/gocircuit/circuit/kit/element/valve"
	"github.com/gocircuit/circuit/use/circuit"
	xerrors "github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XTerminal{})
}

type XTerminal struct {
	t *Terminal
}

func (x XTerminal) Walk(walk []string) circuit.X {
	t := x.t.Walk(walk)
	if t == nil {
		return nil
	}
	return circuit.Ref(XTerminal{t})
}

func (x XTerminal) View() map[string]circuit.X {
	u := make(map[string]circuit.X)
	for p, q := range x.t.View() {
		u[p] = circuit.Ref(XTerminal{q})
	}
	return u
}

func (x XTerminal) Make(kind string, arg interface{}) (xelm circuit.X, err error) {
	elm, err := x.t.Make(kind, arg)
	if err != nil {
		return nil, err
	}
	return elm.X(), nil
}

func (x XTerminal) Get() (string, circuit.X) {
	kind, elm := x.t.Get()
	if elm == nil {
		return "", nil
	}
	return kind, elm.X()
}

func (x XTerminal) Scrub() {
	x.t.Scrub()
}

type YTerminal struct {
	X circuit.X
}

func (y YTerminal) Walk(walk []string) YTerminal {
	return YTerminal{
		y.X.Call("Walk")[0].(circuit.X),
	}
}

func (y YTerminal) View() map[string]YTerminal {
	u := make(map[string]YTerminal)
	for n, x := range y.X.Call("View")[0].(map[string]circuit.X) {
		u[n] = YTerminal{x}
	}
	return u
}

func (y YTerminal) Make(kind string, arg interface{}) (yelm interface{}, err error) {
	r := y.X.Call("Make")
	if err = xerrors.Unpack(r[1]); err != nil {
		return nil, err
	}
	switch kind {
	case Chan:
		return valve.YValve{r[0].(circuit.X)}, nil
	case Proc:
		return proc.YProc{r[0].(circuit.X)}, nil
	}
	return nil, errors.New("element kind not supported")
}

func (y YTerminal) Get() (kind string, yelm interface{}) {
	r := y.X.Call("Get")
	kind = r[0].(string)
	switch kind {
	case Chan:
		return Chan, valve.YValve{r[1].(circuit.X)}
	case Proc:
		return Proc, proc.YProc{r[1].(circuit.X)}
	}
	return "", nil
}

func (y YTerminal) Scrub() {
	y.X.Call("Scrub")
}
