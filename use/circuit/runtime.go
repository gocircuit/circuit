// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package circuit

import (
	"github.com/gocircuit/circuit/use/n"
	"github.com/gocircuit/circuit/use/worker"
)

type runtime interface {
	// Low-level
	WorkerAddr() n.Addr
	SetBoot(interface{})
	Kill(n.Addr) error

	// Spawn mechanism
	Spawn(worker.Host, []string, Func, ...interface{}) ([]interface{}, n.Addr, error)
	RunInBack(func())

	// Cross-services
	Dial(n.Addr, string) PermX
	DialSelf(string) interface{}
	TryDial(n.Addr, string) (PermX, error)
	Listen(string, interface{})

	// Persistence of PermX values
	Export(...interface{}) interface{}
	Import(interface{}) ([]interface{}, string, error)

	// Cross-interfaces
	Ref(interface{}) X
	PermRef(interface{}) PermX

	// Type system
	RegisterValue(interface{})
	RegisterFunc(Func)

	// Utility
	Hang()
}
