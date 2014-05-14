// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"github.com/gocircuit/circuit/kit/anchor"
	"github.com/gocircuit/circuit/kit/element/proc"
	"github.com/gocircuit/circuit/kit/element/valve"
	"github.com/gocircuit/circuit/kit/kinfolk"
)

type Anchor interface {
	Worker() string
	Walk(walk []string) Anchor
	View() map[string]Anchor
	MakeChan(n int) (Chan, error)
	MakeProc(cmd Cmd) (Proc, error)
	Get() interface{}
	Scrub()
}

type terminal struct {
	y anchor.YTerminal
	k kinfolk.KinXID
}

func (t terminal) Worker() string {
	return t.k.ID.String()
}

func (t terminal) Walk(walk []string) Anchor {
	return terminal{ y: t.y.Walk(walk), k: t.k }
}

func (t terminal) View() map[string]Anchor {
	v := t.y.View()
	w := make(map[string]Anchor)
	for name, y := range v {
		w[name] = terminal{ y: y, k: t.k }
	}
	return w
}

func (t terminal) MakeChan(n int) (Chan, error) {
	yvalve, err := t.y.Make(anchor.Chan, n)
	if err != nil {
		return nil, err
	}
	return yvalveChan{yvalve.(valve.YValve)}, nil
}

func (t terminal) MakeProc(cmd Cmd) (Proc, error) {
	yproc, err := t.y.Make(anchor.Proc, cmd.retype())
	if err != nil {
		return nil, err
	}
	return yprocProc{yproc.(proc.YProc)}, nil
}

func (t terminal) Get() interface{} {
	kind, y := t.y.Get()
	if y == nil {
		return nil
	}
	switch kind {
	case anchor.Chan:
		return yvalveChan{y.(valve.YValve)}
	case anchor.Proc:
		return yprocProc{y.(proc.YProc)}
	}
	panic("client/circuit mismatch")
}

func (t terminal) Scrub() {
	t.y.Scrub()
}
