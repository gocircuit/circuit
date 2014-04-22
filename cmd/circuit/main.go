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

	"github.com/gocircuit/circuit/kit/fs/bridge/fuserh"
	"github.com/gocircuit/circuit/kit/kinfolk"
	"github.com/gocircuit/circuit/kit/kinfolk/locus"
	"github.com/gocircuit/circuit/kit/fs/rh"
	nsrv "github.com/gocircuit/circuit/kit/fs/namespace/server"
	"github.com/gocircuit/circuit/kit/fs/client"
	"github.com/gocircuit/circuit/kit/fs/bridge/rhunix"
	"github.com/gocircuit/circuit/kit/shell"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
)

var (
	// System-related
	flagAddr         = flag.String("a", "", "External network address for the worker; required non-empty")
	flagDir          = flag.String("d", "", "Directory to lock; use a new temporary directory, if empty")
	flagWorkerID     = flag.String("i", "", "ID of the worker; choose randomly, if empty")
	flagJoin         = flag.String("j", "", "Join a worker URL. Start a genus worker, if empty")
	flagMount        = flag.String("m", "/"+n.Scheme, "Mount directory within local file system; required non-empty")
	// Resource-related
	flagFS           = flag.String("fs", "", "Local FS directory to share with others; none, if empty")
	// Internal helper roles of the circuit executable
	flagSysOpenRead  = flag.String("sysread", "", "(sys) Open a file for reading.")
	flagSysOpenWrite = flag.String("syswrite", "", "(sys) Open a file for writing.")
)

const (
	KinFolk   = "kin"
	LocusFolk = "locus"
)

func main() {
	flag.Parse()
	if *flagSysOpenRead != "" {
		mainSysOpenRead(*flagSysOpenRead)
		panic(0)
	}
	if *flagSysOpenWrite != "" {
		mainSysOpenWrite(*flagSysOpenWrite)
		panic(0)
	}

	println("CIRCUIT 2014 gocircuit.org")

	// Parse and verify arguments
	if *flagAddr == "" {
		log.Fatal("network address flag -a not given")
	}
	var err error
	// Joining?
	var join n.Addr
	if *flagJoin != "" {
		if join, err = n.ParseAddr(*flagJoin); err != nil {
			log.Fatalf("join address does not parse (%s)", err)
		}
	}
	// Working directory
	if *flagDir == "" {
		*flagDir = path.Join(os.TempDir(), fmt.Sprintf("%s-%%W-P%04d", n.Scheme, os.Getpid()))
	}
	// Mount point
	if *flagMount == "" {
		log.Fatalf("no mount point, flag -m, given")
	}

	// Start circuit runtime
	c := &Config{
		Addr:     *flagAddr,
		WorkerID: *flagWorkerID,
		Dir:      *flagDir,
	}
	load(c)

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Fatalf("kin join (%#v)", r)
	// 	}
	// }()

	// Share local file system?
	var xfs rh.Server
	if *flagFS != "" {
		*flagFS = path.Clean(*flagFS)
		log.Println("sharing local file system", *flagFS)
		if xfs, err = rhunix.New(fmt.Sprintf("%s.xfs·%s", circuit.WorkerAddr().WorkerID(), *flagFS), *flagFS); err != nil {
			log.Fatalf("local file system server start (%s)", err)
		}
	}

	//
	var (
		rsc = locus.Resources{
			TubeTopic: LocusFolk,
			FS:        xfs,
		}
		xkin kinfolk.ExoKin
	)
	var xj circuit.PermX
	dontPanic(
		func() {
			xj = circuit.Dial(join, KinFolk)
		}, 
		"join connect",
	)
	//
	rsc.Kin, xkin, rsc.KinJoin, rsc.KinLeave = kinfolk.NewKin(xj) // Join Kin network
	lcs := locus.NewLocus(&rsc)                                   // Join Locus network

	//
	clientDir := client.NewDir(lcs.ServerDir(), lcs.Peer, lcs.GetPeers)    // Mount focus namespace
	mount, err := fuserh.Mount(*flagMount, nsrv.NewSession(clientDir), 50) // Mount FUSE file system
	if err != nil {
		log.Fatalf("fuse mounting %s (%s)", *flagMount, err)
	}

	//
	circuit.Listen(KinFolk, xkin) // Start kin services
	circuit.Listen("cons", shell.NewXShell(*flagDir))

	// Survive until the local mount is broken.
	log.Println(n.Scheme, "started successfully.")
	if err = mount.EOF(); err != nil {
		log.Fatalf("mount lost (%s)", err)
	}
	log.Println("bye.")
}

func dontPanic(call func(), ifPanic string) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("%s (%s)", ifPanic, r)
		}
	}()
	call()
}
