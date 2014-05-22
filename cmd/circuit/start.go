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
	"os"
	"path"

	"github.com/gocircuit/circuit/kit/kinfolk"
	"github.com/gocircuit/circuit/kit/kinfolk/locus"

	// "github.com/gocircuit/circuit/kit/shell"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"

	"github.com/gocircuit/circuit/github.com/codegangsta/cli"
)

func server(c *cli.Context) {
	println("CIRCUIT 2014 gocircuit.org")
	log.Println("Starting circuit server")
	// parse arguments
	if !c.IsSet("addr") {
		log.Fatal("server network address not given; use -addr")
	}
	var err error
	var join n.Addr
	if c.IsSet("join") {
		if join, err = n.ParseAddr(c.String("join")); err != nil {
			log.Fatalf("join address does not parse (%s)", err)
		}
	}
	var mutexDir string
	if !c.IsSet("mutex") {
		mutexDir = path.Join(os.TempDir(), fmt.Sprintf("%s-%%W-P%04d", n.Scheme, os.Getpid()))
	} else {
		mutexDir = c.String("mutex")
	}

	// start circuit runtime
	load(c.String("addr"), mutexDir, readkey(c))

	// kinfolk join
	var xjoin circuit.PermX
	dontPanic(func() { 
		xjoin = circuit.Dial(join, KinfolkName) 
	}, "join")

	// locus
	kin, xkin, kinJoin, kinLeave := kinfolk.NewKin(xjoin)
	xlocus := locus.NewLocus(kin, kinJoin, kinLeave)

	circuit.Listen(KinfolkName, xkin)
	circuit.Listen(LocusName, xlocus)

	<-(chan int)(nil)
}

const KinfolkName = "kin"
const LocusName = "locus"

func dontPanic(call func(), ifPanic string) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("%s (%s)", ifPanic, r)
		}
	}()
	call()
}
