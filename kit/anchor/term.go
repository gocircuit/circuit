// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"errors"
	"io"
	"log"

	"github.com/gocircuit/circuit/kit/element/proc"
	"github.com/gocircuit/circuit/kit/element/valve"
	"github.com/gocircuit/circuit/use/circuit"
)

type Element interface {
	Scrub()
	X() circuit.X
}

const (
	Chan = "chan"
	Proc = "proc"
	//Pipe = "pipe"
	//Mutex = "mutex"
)

// Terminal presents a facade to *Anchor with added element manipulation methods
type Terminal Anchor

func (t *Terminal) carrier() *Anchor {
	return (*Anchor)(t)
}

func (t *Terminal) Walk(walk []string) *Terminal {
	return (*Terminal)(t.carrier().Walk(walk))
}

func (t *Terminal) View() map[string]*Terminal {
	r := make(map[string]*Terminal)
	for n, a := range t.carrier().View() {
		r[n] = (*Terminal)(a)
	}
	return r
}

type urn struct {
	kind string
	elem Element // valve.Valve, proc.Proc, etc
}

func (t *Terminal) Make(kind string, arg interface{}) (elem Element, err error) {
	log.Printf("%s make %s with %v", t.carrier().Path(), kind, arg)
	t.carrier().TxLock()
	defer t.carrier().TxUnlock()
	if t.carrier().Get() != nil {
		return nil, errors.New("anchor already has an element")
	}
	switch kind {
	case Chan:
		capacity, ok := arg.(int)
		if !ok {
			return nil, errors.New("invalid argument")
		}
		u := &urn{
			kind: Chan,
			elem: &scrubValve{t, valve.MakeValve(capacity)},
		}
		t.carrier().Set(u)
		return u.elem, nil
	case Proc:
		cmd, ok := arg.(proc.Cmd)
		if !ok {
			return nil, errors.New("invalid argument")
		}
		u := &urn{
			kind: Proc,
			elem: proc.MakeProc(cmd),
		}
		t.carrier().Set(u)
		go func() {
			defer func() {
				recover()
			}()
			defer t.Scrub()
			u.elem.(proc.Proc).Wait()
		}()
		return u.elem, nil
	}
	return nil, errors.New("element kind not known")
}

func (t *Terminal) Get() (string, Element) {
	t.carrier().TxLock()
	defer t.carrier().TxUnlock()
	v := t.carrier().Get()
	if v == nil {
		return "", nil
	}
	return v.(*urn).kind, v.(*urn).elem
}

func (t *Terminal) Scrub() {
	log.Printf("scrubbing %s", t.carrier().Path())
	t.carrier().TxLock()
	defer t.carrier().TxUnlock()
	u, ok := t.carrier().Get().(*urn)
	if !ok {
		return
	}
	u.elem.Scrub()
	t.carrier().Set(nil)
}

type scrubValve struct {
	t *Terminal
	valve.Valve
}

func (v *scrubValve) Close() error {
	defer func() {
		if v.Valve.IsDone() {
			v.t.Scrub()
		}
	}()
	return v.Valve.Close()
}

func (v *scrubValve) Recv() (io.ReadCloser, error) {
	defer func() {
		if v.Valve.IsDone() {
			v.t.Scrub()
		}
	}()
	return v.Valve.Recv()
}
