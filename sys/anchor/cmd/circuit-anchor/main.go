// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package circuit-anchor is an Anchor/Graveyard File System server
package main

import (
	"flag"
	"log"

	"github.com/gocircuit/circuit/sys/anchor"
	"github.com/gocircuit/circuit/sys/anchor/xy"
	"github.com/gocircuit/circuit/use/circuit"
)

var (
	flagAddr = flag.String("addr", "127.0.0.1:49001", "Bind address of the AGFS server")
	flagID   = flag.String("id", "", "Specify worker ID")
)

func main() {
	flag.Parse()
	boot(*flagAddr, *flagID)
	xsys := (*xy.XSys)(anchor.NewSystem())
	circuit.Listen("anchor", xsys)
	log.Println("∞ anchor/graveyard file system started successfully.")
	circuit.Hang()
}
