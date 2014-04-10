// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package dash

import (
	"os"
	"testing"

	"github.com/gocircuit/circuit/kit/debug"
	_ "github.com/gocircuit/circuit/kit/debug/ctrlc"
	"github.com/gocircuit/circuit/kit/fs/fuse"
	"github.com/gocircuit/circuit/kit/fs/bridge/fuserh"
	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
)

var testMount = "/t"

func init() {
	debug.OnSignal(func(os.Signal) {
		fuse.Umount(testMount)
	})
	os.MkdirAll(testMount, 0777)
}

// TestMount
func TestMount(t *testing.T) {
	slash := NewDir("/") // dash
	ssn, err := dir.NewServer(slash).SignIn("user", "")
	if err != nil {
		panic(err)
	}
	// Start FUSE, serving a namespace
	if testMount == "" {
		panic("need env GOTEST_MOUNT")
	}
	fr, err := fuserh.Mount(testMount, ssn, 5)
	if err != nil {
		panic(err)
	}
	defer fuse.Umount(testMount)

	// Wait until end of serving
	println("serving...")
	if err := fr.EOF(); err != nil {
		println("EOF with", err.Error())
	}
	println("EOF ok")
}
