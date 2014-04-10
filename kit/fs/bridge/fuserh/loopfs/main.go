// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"fmt"
	"os"

	_ "github.com/gocircuit/circuit/kit/debug/ctrlc" // Ctrl-C panics the process
	"github.com/gocircuit/circuit/kit/fs/bridge/fuserh"
	"github.com/gocircuit/circuit/kit/fs/rh/rhunix"
)

func main() {
	if len(os.Args) != 3 {
		println("usage: loopfs source_dir mount_dir")
		os.Exit(1)
	}
	ru, err := rhunix.New("loopfs", os.Args[1])
	if err != nil {
		panic(err)
	}
	ssn, err := ru.SignIn("1", "/")
	if err != nil {
		panic(err)
	}
	fr, err := fuserh.Mount(os.Args[2], ssn, 5)
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("loopfs serving %s on %s ...", os.Args[1], os.Args[2]))
	if err := fr.EOF(); err != nil {
		panic(err)
	}
}
