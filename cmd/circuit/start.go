// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// This package provides the executable program for the resource-sharing circuit app
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"

	"github.com/gocircuit/circuit/kit/discover"
	"github.com/gocircuit/circuit/kinfolk"
	"github.com/gocircuit/circuit/kinfolk/locus"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"

	"github.com/gocircuit/circuit/github.com/codegangsta/cli"
)

func server(c *cli.Context) {
	println("CIRCUIT 2014 gocircuit.org")

	// parse arguments
	var err error
	// network address for server to bind
	if !c.IsSet("addr") {
		log.Fatal("server network address not given; use -addr")
	}
	// join address of another circuit server
	var join n.Addr
	if c.IsSet("join") {
		if join, err = n.ParseAddr(c.String("join")); err != nil {
			log.Fatalf("join address does not parse (%s)", err)
		}
	}
	// discover system udp multicast address
	var disc *net.UDPAddr
	if c.IsSet("discover") {
		if disc, err = net.ResolveUDPAddr("udp", c.String("discover")); err != nil {
			log.Fatalf("udp multicast address for discovery does not parse (%s)", err)
		}
	}
	// server instance working directory
	var varDir string
	if !c.IsSet("var") {
		varDir = path.Join(os.TempDir(), fmt.Sprintf("%s-%%W-P%04d", n.Scheme, os.Getpid()))
	} else {
		varDir = c.String("var")
	}

	// start circuit runtime
	addr := load(c.String("addr"), varDir, readkey(c))

	// kinfolk + locus
	kin, xkin, kinJoin, kinLeave := kinfolk.NewKin()
	xlocus := locus.NewLocus(kin, kinJoin, kinLeave)

	// joining
	switch {
	case join != nil:
		kin.ReJoin(join)
	case disc != nil:
		_, ch := discover.New(disc, []byte(addr.String()))
		go func() {
			for ja := range ch {
				join, err := n.ParseAddr(string(ja))
				if err != nil {
					log.Printf("Unrecognized discovery packets (%v)", err)
					continue // skip messages that don't parse
				}
				kin.ReJoin(join)
			}
		}()
	default:
		log.Println("Singleton server.")
	}

	circuit.Listen(kinfolk.ServiceName, xkin)
	circuit.Listen(LocusName, xlocus)

	<-(chan int)(nil)
}

const LocusName = "locus"

func dontPanic(call func(), ifPanic string) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("%s (%s)", ifPanic, r)
		}
	}()
	call()
}
