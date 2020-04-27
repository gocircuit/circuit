// +build !windows
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

	_ "github.com/hoijui/circuit/kit/debug/kill"
	"github.com/hoijui/circuit/kit/lockfile"
	"github.com/hoijui/circuit/sys/lang"
	_ "github.com/hoijui/circuit/sys/tele"
	"github.com/hoijui/circuit/use/circuit"
	"github.com/hoijui/circuit/use/n"
)

func load(addr *net.TCPAddr, vardir string, key []byte) n.Addr {
	//debug.InstallCtrlCPanic()

	// Randomize execution
	rand.Seed(time.Now().UnixNano())

	// Generate worker ID
	id := n.ChooseWorkerID()
	vardir = strings.Replace(vardir, "%W", id.String(), 1)

	// Ensure chroot directory exists and we have access to it
	dir, err := filepath.Abs(vardir)
	if err != nil {
		log.Fatalf("abs (%s)", err)
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalf("mkdir %s (%s)", dir, err)
	}

	// Create a lock file in the chroot directory so its not managed by two circuit instances at the same time
	lockname := path.Join(dir, ".lock")
	if _, err := lockfile.Create(lockname); err != nil {
		log.Fatalf("obtain lock (%s)\n", err)
	}
	log.Printf("Created and locked %s", lockname)

	// Initialize networking
	if len(key) > 0 {
		log.Println("Using symmetric HMAC authentication and RC4 encryption.")
	}
	t := n.NewTransport(id, addr, key)
	fmt.Println(t.Addr().String())

	// Initialize language runtime
	circuit.Bind(lang.New(t))
	return t.Addr()
}
