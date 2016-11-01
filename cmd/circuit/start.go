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

	"github.com/gocircuit/circuit/element/docker"
	"github.com/gocircuit/circuit/kit/assemble"
	"github.com/gocircuit/circuit/tissue"
	"github.com/gocircuit/circuit/tissue/locus"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
	"github.com/pkg/errors"

	"github.com/urfave/cli"
)

func server(c *cli.Context) (err error) {
	println("CIRCUIT 2015 gocircuit.org")

	if c.Bool("docker") {
		cmd, e := docker.Init()
		if e != nil {
			return errors.Wrapf(e, "cannot use docker: %v", e)
		}
		log.Printf("Enabling docker elements, using %s", cmd)
	}
	// parse arguments
	var tcpaddr = parseAddr(c) // server bind address
	var join n.Addr            // join address of another circuit server
	if c.IsSet("join") {
		if join, err = n.ParseAddr(c.String("join")); err != nil {
			return errors.Wrapf(err, "join address does not parse (%s)", err)
		}
	}
	var multicast = parseDiscover(c)
	// server instance working directory
	var varDir string
	if !c.IsSet("var") {
		varDir = path.Join(os.TempDir(), fmt.Sprintf("%s-%%W-P%04d", n.Scheme, os.Getpid()))
	} else {
		varDir = c.String("var")
	}

	// start circuit runtime
	addr := load(tcpaddr, varDir, readkey(c))

	// tissue + locus
	kin, xkin, rip := tissue.NewKin()
	xlocus := locus.NewLocus(kin, rip)

	// joining
	switch {
	case join != nil:
		kin.ReJoin(join)
	case multicast != nil:
		log.Printf("Using UDP multicast discovery on address %s", multicast.String())
		go assemble.NewAssembler(addr, multicast).AssembleServer(
			func(joinAddr n.Addr) {
				kin.ReJoin(joinAddr)
			},
		)
	default:
		log.Println("Singleton server.")
	}

	circuit.Listen(tissue.ServiceName, xkin)
	circuit.Listen(LocusName, xlocus)

	<-(chan int)(nil)
	return nil
}

func parseDiscover(c *cli.Context) *net.UDPAddr {
	src := c.String("discover")
	if src == "" {
		return nil
	}
	multicast, err := net.ResolveUDPAddr("udp", src)
	if err != nil {
		log.Fatalf("udp multicast address for discovery and assembly does not parse (%s)", err)
	}
	return multicast
}

func parseAddr(c *cli.Context) *net.TCPAddr {
	switch {
	case c.String("addr") != "":
		addr, err := net.ResolveTCPAddr("tcp", c.String("addr"))
		if err != nil {
			log.Fatalf("resolve %s (%s)\n", addr, err)
		}
		if len(addr.IP) == 0 {
			addr.IP = net.IPv4zero
		}
		return addr

	case c.String("if") != "":
		ifc, err := net.InterfaceByName(c.String("if"))
		if err != nil {
			log.Fatalf("interface %s not found (%v)", c.String("if"), err)
		}
		addrs, err := ifc.Addrs()
		if err != nil {
			log.Fatalf("interface address cannot be retrieved (%v)", err)
		}
		if len(addrs) == 0 {
			log.Fatalf("no addresses associated with this interface")
		}
		for _, a := range addrs { // pick the IPv4 one
			ipn := a.(*net.IPNet)
			if ipn.IP.To4() == nil {
				continue
			}
			return &net.TCPAddr{IP: ipn.IP}
		}
		log.Fatal("specified interface has no IPv4 addresses")
	default:
		log.Fatal("either an -addr or an -if option is required to start a server")
	}
	panic(0)
}

// LocusName ???
const LocusName = "locus"

func dontPanic(call func(), ifPanic string) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("%s (%s)", ifPanic, r)
		}
	}()
	call()
}
