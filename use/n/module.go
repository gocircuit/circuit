// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package n

import (
	"github.com/hoijui/circuit/kit/module"
	"net"
)

var mod = module.Slot{Name: "network"}

func Bind(v System) {
	mod.Set(v)
}

func get() System {
	return mod.Get().(System)
}

var workeraddr Addr

func ServerAddr() Addr {
	return workeraddr
}

// NewTransport creates a new transport framework for the given local address.
func NewTransport(workerID WorkerID, addr net.Addr, key []byte) Transport {
	t := get().NewTransport(workerID, addr, key)
	workeraddr = t.Addr()
	return t
}

func ParseNetAddr(s string) (net.Addr, error) {
	return get().ParseNetAddr(s)
}

func ParseAddr(s string) (Addr, error) {
	return get().ParseAddr(s)
}
