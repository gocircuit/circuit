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
)

type Element interface {
	Scrub()
}

type ElementAnchor Anchor

func (ea *ElementAnchor) carrier() *Anchor {
	return (*Anchor)(ea)
}

func (ea *ElementAnchor) Walk(walk []string) *ElementAnchor {
	return (*ElementAnchor)(ea.carrier().Walk(walk))
}

func (ea *ElementAnchor) View() map[string]struct{} {
	return ea.carrier().View()
}

type urn struct {
	kind string
	elem Element // chan, proc, etc
}

func (ea *ElementAnchor) Make(kind string, arg interface{}) (elem Element, err error) {
	ea.carrier().TxLock()
	defer ea.carrier().TxUnlock()
	if ea.carrier().Get() != nil {
		return nil, errors.New("anchor already has an element")
	}
	switch kind {
	case "chan":
		capacity, ok := arg.(int)
		if !ok {
			return nil, errors.New("invalid argument")
		}
		u := &urn{
			kind: "chan",
			elem: valve.MakeValve(capacity),
		}
		ea.carrier().Set(u)
		return u.elem, nil
	case "proc":
		cmd, ok := arg.(*proc.Cmd)
		if !ok {
			return nil, errors.New("invalid argument")
		}
		u := &urn{
			kind: "proc",
			elem: proc.MakeProc(cmd),
		}
		ea.carrier().Set(u)
		return u.elem, nil
	}
	return nil, errors.New("element kind not known")
}

func (ea *ElementAnchor) Get(kind string, arg interface{}) (string, Element) {
	ea.carrier().TxLock()
	defer ea.carrier().TxUnlock()
	v := ea.carrier().Get()
	if v == nil {
		return "", nil
	}
	return v.(*urn).kind, v.(*urn).elem
}

// ??
func (ea *ElementAnchor) Scrub() {
	ea.carrier().TxLock()
	defer ea.carrier().TxUnlock()
	u, ok := ea.carrier().Get().(*urn)
	if !ok {
		return
	}
	u.elem.Scrub()
	ea.carrier().Set(nil)
}
