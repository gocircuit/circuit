// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// This package provides the executable program for the resource-sharing circuit app
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/gocircuit/circuit/kit/kinfolk"
	"github.com/gocircuit/circuit/kit/kinfolk/locus"

	//"github.com/gocircuit/circuit/kit/shell"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
)

var (
	flagAddr         = flag.String("addr", "", "Network address to use")
	flagDir          = flag.String("lock", "", "Directory to lock to prevent duplication")
	flagJoin         = flag.String("join", "", "Join an existing network of circuit workers")
)

func main() {
	flag.Parse()
	println("CIRCUIT 2014 gocircuit.org")

	// parse arguments
	if *flagAddr == "" {
		log.Fatal("network address not given")
	}
	var err error
	var join n.Addr
	if *flagJoin != "" {
		if join, err = n.ParseAddr(*flagJoin); err != nil {
			log.Fatalf("join address does not parse (%s)", err)
		}
	}
	if *flagDir == "" {
		*flagDir = path.Join(os.TempDir(), fmt.Sprintf("%s-%%W-P%04d", n.Scheme, os.Getpid()))
	}

	// start circuit runtime
	c := &Config{Addr: *flagAddr, Dir: *flagDir}
	load(c)

	// kinfolk join
	var xjoin circuit.PermX
	dontPanic(func() { 
		xjoin = circuit.Dial(join, KinfolkName) 
	}, "join")

	// locus
	kin, xkin, kinJoin, kinLeave := kinfolk.NewKin(xjoin)
	/*l := */ locus.NewLocus(kin, kinJoin, kinLeave)

	circuit.Listen(KinfolkName, xkin) // Start kin services
	//circuit.Listen(EyeName, ???)
	//circuit.Listen("cons", shell.NewXShell(*flagDir))

	<-(chan int)(nil)
}

const KinfolkName = "kin"
const EyeName = "eye"

func dontPanic(call func(), ifPanic string) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("%s (%s)", ifPanic, r)
		}
	}()
	call()
}
