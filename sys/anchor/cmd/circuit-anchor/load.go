// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	_ "github.com/gocircuit/circuit/kit/debug/ctrlc"
	_ "github.com/gocircuit/circuit/kit/debug/kill"

	// Sys
	"github.com/gocircuit/circuit/sys/lang"
	_ "github.com/gocircuit/circuit/sys/tele"

	// Use
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
)

func boot(addr, workerID string) {

	// Randomize execution
	rand.Seed(time.Now().UnixNano())

	// Generate worker ID
	var id n.WorkerID
	if workerID == "" {
		id = n.ChooseWorkerID()
	} else {
		id = n.ParseOrHashWorkerID(workerID)
	}
	//fmt.Println("WorkerID:", id)

	// Initialize networking
	bindaddr_, err := n.ParseNetAddr(addr)
	if err != nil {
		log.Printf("resolve %s (%s)\n", addr, err)
		os.Exit(1)
	}
	bindaddr := bindaddr_.(*net.TCPAddr)
	if len(bindaddr.IP) == 0 {
		bindaddr.IP = net.IPv4zero
	}
	t := n.NewTransport(id, bindaddr)

	// Print port number
	// port := t.Addr().NetAddr().(*net.TCPAddr).Port
	// fmt.Println("CircuitPort:", port)
	fmt.Println("Addr:", t.Addr().String())

	// Initialize language runtime
	circuit.Bind(lang.New(t))
}
