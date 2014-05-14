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
	"path"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/gocircuit/circuit/kit/debug/ctrlc"
	_ "github.com/gocircuit/circuit/kit/debug/kill"
	"github.com/gocircuit/circuit/kit/lockfile"

	// Sys
	"github.com/gocircuit/circuit/sys/lang"
	_ "github.com/gocircuit/circuit/sys/tele"

	// Use
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
)

func load(addr, mutex string) {

	// Randomize execution
	rand.Seed(time.Now().UnixNano())

	// Generate worker ID
	id := n.ChooseWorkerID()
	mutex = strings.Replace(mutex, "%W", id.String(), 1)

	// Ensure chroot directory exists and we have access to it
	dir, err := filepath.Abs(mutex)
	if err != nil {
		log.Fatalf("abs (%s)", err)
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalf("mkdir %s (%s)", dir, err)
	}

	// Create a lock file in the chroot directory so its not managed by two REX instances at the same time
	lockname := path.Join(dir, ".lock")
	if _, err := lockfile.Create(lockname); err != nil {
		log.Fatalf("obtain lock (%s)\n", err)
	}
	log.Printf("Created and locked %s", lockname)

	// Initialize networking
	bindaddr_, err := n.ParseNetAddr(addr)
	if err != nil {
		log.Fatalf("resolve %s (%s)\n", addr, err)
	}
	bindaddr := bindaddr_.(*net.TCPAddr)
	if len(bindaddr.IP) == 0 {
		bindaddr.IP = net.IPv4zero
	}
	t := n.NewTransport(id, bindaddr)
	fmt.Println(t.Addr().String())

	// Initialize language runtime
	circuit.Bind(lang.New(t))
}
