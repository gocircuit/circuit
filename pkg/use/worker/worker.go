// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package worker is a facade for the circuit spawning mechanism module
package worker

import (
	"github.com/hoijui/circuit/pkg/kit/module"
	"github.com/hoijui/circuit/pkg/use/n"
)

// Host is a location where a new worker can be executed.
type Host interface{}

// Spawn starts a new worker process on host and registers it under the given
// anchors directories in the anchor file system. On success, Spawn returns
// the address of the new work. Spawn is a low-level function. The spawned
// worker will wait idle for further interaction. It is the caller's responsibility
// to manage the lifespan of the newworker.
func Spawn(host Host, anchors ...string) (n.Addr, error) {
	return get().Spawn(host, anchors...)
}

// Kill kills the circuit worker with the given addr
func Kill(addr n.Addr) error {
	return get().Kill(addr)
}

type commander interface {
	Spawn(Host, ...string) (n.Addr, error)
	Kill(n.Addr) error
}

// Binding mechanism
var mod = module.Slot{Name: "worker"}

// Bind is used internally to bind an implementation of this package to the public methods of this package
func Bind(v commander) {
	mod.Set(v)
}

func get() commander {
	return mod.Get().(commander)
}
