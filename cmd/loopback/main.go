// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"flag"
	"os"

	_ "github.com/gocircuit/circuit/kit/debug/ctrlc"
	"github.com/gocircuit/circuit/kit/fs/element/knot"
	"github.com/gocircuit/circuit/kit/fs/fuse"
	"github.com/gocircuit/circuit/kit/fs/helper"
	"github.com/gocircuit/circuit/kit/fs/bridge/fuserh"
	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
)

var flagMount = flag.String("m", "", "mount point")

func main() {
	flag.Parse()
	helper.Main()

	slash := knot.NewDir("/")
	ssn, err := dir.NewServer(slash).SignIn("", "")
	if err != nil {
		panic(err)
	}

	// prep for mounting
	if *flagMount == "" {
		panic("need mount point -m")
	}
	os.MkdirAll(*flagMount, 0777)

	// mount
	fr, err := fuserh.Mount(*flagMount, ssn, 5)
	if err != nil {
		panic(err)
	}
	defer fuse.Umount(*flagMount)

	// wait for end
	println("serving ...")
	if err := fr.EOF(); err != nil {
		println("end of connection with", err.Error())
	}
	println("disconnected.")
}
