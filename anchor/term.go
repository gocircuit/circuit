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

	ds "github.com/gocircuit/circuit/client/docker"
	"github.com/gocircuit/circuit/element/dns"
	"github.com/gocircuit/circuit/element/docker"
	"github.com/gocircuit/circuit/element/proc"
	srv "github.com/gocircuit/circuit/element/server"
	"github.com/gocircuit/circuit/element/valve"
	"github.com/gocircuit/circuit/kit/pubsub"
	"github.com/gocircuit/circuit/use/circuit"
)

type Element interface {
	Scrub()
	X() circuit.X
}

const (
	Server     = "server"
	Chan       = "chan"
	Proc       = "proc"
	Docker     = "docker"
	Nameserver = "dns"
	OnJoin     = "@join"
	OnLeave    = "@leave"
)

// Terminal presents a facade to *Anchor with added element manipulation methods
type Terminal struct {
	genus  Genus
	anchor *Anchor
}

type Genus interface {
	NewArrivals() pubsub.Consumer
	NewDepartures() pubsub.Consumer
}

// NewTerm create the root node of a new anchor file system.
func NewTerm(name string, genus Genus) (*Terminal, circuit.PermX) {
	t := &Terminal{
		genus:  genus,
		anchor: newAnchor(nil, name).use(),
	}
	return t, circuit.PermRef(XTerminal{t})
}

func (t *Terminal) carrier() *Anchor {
	return t.anchor
}

func (t *Terminal) Walk(walk []string) *Terminal {
	return &Terminal{
		genus:  t.genus,
		anchor: t.carrier().Walk(walk),
	}
}

func (t *Terminal) Path() string {
	return t.carrier().Path()
}

func (t *Terminal) View() map[string]*Terminal {
	r := make(map[string]*Terminal)
	for n, a := range t.carrier().View() {
		r[n] = &Terminal{
			genus:  t.genus,
			anchor: a,
		}
	}
	return r
}

type urn struct {
	kind string
	elem Element // valve.Valve, proc.Proc, etc
}

func (t *Terminal) Attach(kind string, elm Element) {
	if kind != Server {
		panic(0)
	}
	log.Printf("Attaching %s as %s", t.carrier().Path(), kind)
	t.carrier().TxLock()
	defer t.carrier().TxUnlock()
	if t.carrier().Get() != nil {
		panic(0)
	}
	u := &urn{
		kind: kind,
		elem: elm,
	}
	t.carrier().Set(u)
}

func (t *Terminal) Make(kind string, arg interface{}) (elem Element, err error) {
	log.Printf("Making %s at %s, using %v", kind, t.carrier().Path(), arg)
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
			if cmd.Scrub {
				defer t.Scrub()
			}
			u.elem.(proc.Proc).Wait()
		}()
		return u.elem, nil

	case Docker:
		run, ok := arg.(ds.Run)
		if !ok {
			return nil, errors.New("invalid argument")
		}
		x, err := docker.MakeContainer(run)
		if err != nil {
			return nil, err
		}
		u := &urn{
			kind: Docker,
			elem: x,
		}
		t.carrier().Set(u)
		go func() {
			defer func() {
				recover()
			}()
			if run.Scrub {
				defer t.Scrub()
			}
			u.elem.(docker.Container).Wait()
		}()
		return u.elem, nil

	case Nameserver:
		ns, err := dns.MakeNameserver(arg.(string))
		if err != nil {
			return nil, err
		}
		u := &urn{
			kind: Nameserver,
			elem: ns,
		}
		t.carrier().Set(u)
		return u.elem, nil

	case OnJoin:
		u := &urn{
			kind: OnJoin,
			elem: t.genus.NewArrivals(),
		}
		t.carrier().Set(u)
		return u.elem, nil

	case OnLeave:
		u := &urn{
			kind: OnLeave,
			elem: t.genus.NewDepartures(),
		}
		t.carrier().Set(u)
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
	log.Printf("Scrubbing %s", t.carrier().Path())
	t.carrier().TxLock()
	defer t.carrier().TxUnlock()
	u, ok := t.carrier().Get().(*urn)
	if !ok {
		return
	}
	if _, ok := u.elem.(srv.Server); ok {
		return // Cannot scrub server anchors
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
