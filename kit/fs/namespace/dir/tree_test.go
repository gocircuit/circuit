// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package dir

import (
	"os"
	"testing"

	"github.com/gocircuit/circuit/kit/debug"
	_ "github.com/gocircuit/circuit/kit/debug/ctrlc"
	"github.com/gocircuit/circuit/kit/fs/fuse"
	"github.com/gocircuit/circuit/kit/fs/bridge/fuserh"
	"github.com/gocircuit/circuit/kit/fs/rh/ns"
)

const testMount = "/tmp/nonce"

func init() {
	debug.OnSignal(func(os.Signal) {
		fuse.Umount(testMount)
	})
	os.MkdirAll(testMount, 0777)
}

// TestDir
func TestDir(t *testing.T) {
	slash := NewDir()
	namspac, err := ns.New(slash.FID()).SignIn("user", "")
	if err != nil {
		panic(err)
	}
	// Start FUSE, serving a namespace
	fr, err := fuserh.Mount(testMount, namspac, 5)
	if err != nil {
		panic(err)
	}
	defer fuse.Umount(testMount)

	// Add children
	slash.AddChild("bobo", NewDir().FID())

	// Wait until end of serving
	println("serving...")
	if err := fr.EOF(); err != nil {
		println("EOF with", err.Error())
	}
	println("EOF ok")
}
