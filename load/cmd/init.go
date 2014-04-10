// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package cmd has the side effect of linking the circuit runtime into the importing application
package cmd

import (
	"net"
	"os"
	"path"
	"strconv"
	"strings"

	_ "github.com/gocircuit/circuit/kit/debug/ctrlc"
	_ "github.com/gocircuit/circuit/kit/debug/kill"
	"github.com/gocircuit/circuit/kit/lockfile"

	anchorback "github.com/gocircuit/circuit/sys/anchor/load"
	"github.com/gocircuit/circuit/sys/lang"
	_ "github.com/gocircuit/circuit/sys/tele"
	workerback "github.com/gocircuit/circuit/sys/worker"

	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
	"github.com/gocircuit/circuit/use/worker"

	"github.com/gocircuit/circuit/load/config" // Side-effect of reading in configurations
	"log"
)

func init() {
	log.Println("∞ loading circuit runtime ...")
	defer log.Println("∞ loaded circuit runtime successfully.")

	switch config.Role {
	case config.Main:
		start(false, config.Config.Anchor, config.Config.Deploy, config.Config.Spark)
	case config.Worker:
		start(true, config.Config.Anchor, config.Config.Deploy, config.Config.Spark)
	case config.Daemonizer:
		workerback.Daemonize(config.Config)
	default:
		log.Println("Circuit role unrecognized:", config.Role)
		os.Exit(1)
	}
}

func start(isWorker bool, acfg *config.AnchorConfig, i *config.InstallConfig, s *config.SparkConfig) {
	// If this is a worker, create a lock file in its working directory
	if isWorker {
		if _, err := lockfile.Create("lock"); err != nil {
			log.Printf("Worker cannot obtain lock (%s)\n", err)
			os.Exit(1)
		}
		wd, err := os.Getwd()
		if err != nil {
			log.Println("Error creating lock:", err.Error())
			os.Exit(1)
		}
		log.Println("Created lock in", wd)
	}

	// Initialize the networking module
	worker.Bind(workerback.New(i.LibPath, path.Join(i.BinDir(), i.Worker), i.JailDir()))

	// Initialize transport module
	bindaddr_, err := n.ParseNetAddr(s.BindAddr)
	if err != nil {
		panic(err.Error())
	}
	bindaddr := bindaddr_.(*net.TCPAddr)
	if len(bindaddr.IP) == 0 {
		bindaddr.IP = net.IPv4zero
	}

	t := n.NewTransport(s.ID, bindaddr) // s.Host is not used any more.
	port := t.Addr().NetAddr().(*net.TCPAddr).Port

	// Initialize language runtime
	circuit.Bind(lang.New(t))

	// Connect to anchor file system
	anchorSrvAddr, err := n.ParseAddr(acfg.Addr)
	if err != nil {
		panic(err)
	}
	for i, a := range s.Anchor {
		s.Anchor[i] = path.Clean(strings.TrimSpace(a))
	}
	s.Anchor = append(s.Anchor, "/addr")
	anchorback.Load(anchorSrvAddr, t.Addr(), s.Anchor)

	// Handy printouts
	// fmt.Println("CircuitAddr:", t.Addr().String())

	if isWorker {
		// A worker sends back its PID and runtime port to its invoker (the daemonizer)
		backpipe := os.NewFile(3, "backpipe")
		if _, err := backpipe.WriteString(strconv.Itoa(os.Getpid()) + "\n"); err != nil {
			panic(err)
		}
		if _, err := backpipe.WriteString(strconv.Itoa(port) + "\n"); err != nil {
			panic(err)
		}
		if err := backpipe.Close(); err != nil {
			panic(err)
		}
	}
}
