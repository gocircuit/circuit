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


type Terminal struct {
	y anchor.YTerminal
	k kinfolk.KinXID
}

func (t Terminal) Worker() string {
	return t.k.ID.String()
}

func (t Terminal) Walk(walk []string) Terminal {
	return Terminal{ y: t.y.Walk(walk), k: t.k }
}

func (t Terminal) View() map[string]Terminal {
	v := t.y.View()
	w := make(map[string]Terminal)
	for name, y := range v {
		w[name] = Terminal{ y: y, k: t.k }
	}
	return w
}

func (t Terminal) MakeChan(n int) (Chan, error) {
	yvalve, err := t.y.Make(anchor.Chan, n)
	if err != nil {
		return nil, err
	}
	return yvalveChan{yvalve.(valve.YValve)}, nil
}

func (t Terminal) MakeProc(cmd Cmd) (Proc, error) {
	yproc, err := t.y.Make(anchor.Proc, cmd.retype())
	if err != nil {
		return nil, err
	}
	return yprocProc{yproc.(proc.YProc)}, nil
}

func (t Terminal) Get() interface{} {
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
