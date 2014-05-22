// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"github.com/gocircuit/circuit/kit/anchor"
	"github.com/gocircuit/circuit/element/proc"
	"github.com/gocircuit/circuit/element/valve"
	"github.com/gocircuit/circuit/kit/kinfolk"
)

// An Anchor represents a location in the global anchor namespace of a circuit
// cluster. Anchors are named locations where the user can store and operate
// control primitives, called circuit elements. The anchor namespace hierarchy
// is represented in paths of the form
//
//	/X8817c114d4941522/hello/dolly
//
// The root anchor "/" represents the cluster abstractly and is the only
// anchor within which one cannot create elements or freely-named subanchors.
// The root anchor contains a dynamically changing set of sub-anchors that
// correspond to the live circuit servers in the cluster.
//
// Every anchor, other than "/", can be used to make, store and operate a
// circuit element (a process or a channel). Anchors are created on access, if
// not present, and are garbage-collected when not used or referenced.
// Therefore the interface allows users to access arbitrary paths without
// having to create them first.
//
type Anchor interface {

	// Addr returns the address of the circuit server hosting this anchor.
	Addr() string

	// ServerID returns the ID of the circuit server hosting this anchor.
	// The returned string will look like "X123..."
	ServerID() string

	// Walk traverses the anchor namespace, starting from this anchor along the path in walk.
	// Errors in communication or a missing circuit server condition are reported via panics.
	Walk(walk []string) Anchor

	// View returns the set of its sub-anchors.
	View() map[string]Anchor

	// MakeChan creates a new circuit channel element at this anchor with a given capacity n.
	// If the anchor already stores an element, a non-nil error is returned.
	// Panics indicate that the server hosting the anchor is gone.
	MakeChan(n int) (Chan, error)

	// MakeProc issues the execution of an OS process, described by cmd, at the server hosting the anchor
	// and creates a corresponding circuit process element at this anchor.
	// If the anchor already stores an element, a non-nil error is returned.
	// Panics indicate that the server hosting the anchor is gone.
	MakeProc(cmd Cmd) (Proc, error)

	// Get returns a handle for the circuit element (one of Chan or Proc) stored at this anchor, and nil otherwise. 
	// Panics indicate that the server hosting the anchor and its element has already died.
	Get() interface{}

	// Scrub aborts and abandons the circuit element stored at this anchor, if one is present.
	// If the hosting server is dead, a panic will be issued.
	Scrub()
}

// Split breaks up an anchor path into components.
func Split(walk string) (r []string) {
	var j int
	for i, c := range walk {
		if c != '/' {
			continue
		}
		if i - j > 0 {
			r = append(r, walk[j:i])
		}
		j = i+1
	}
	if len(walk) - j > 0 {
		r = append(r, walk[j:])
	}
	return
}

type terminal struct {
	y anchor.YTerminal
	k kinfolk.KinXID
}

func (t terminal) Addr() string {
	return t.k.X.Addr().String()
}

func (t terminal) ServerID() string {
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
