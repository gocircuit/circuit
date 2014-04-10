// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"flag"
	"math/rand"
	"time"

	//_ "github.com/gocircuit/circuit/kit/debug/ctrlc"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	flagServer = flag.Bool("server", false, "Run in server regime")
	flagClient = flag.Bool("client", false, "Run in client regime")
	flagIn     = flag.String("in", "", "Input address")
	flagOut    = flag.String("out", "", "Output address")
	//flagMax    = flag.Int("max", 100, "Maximum number of concurrent connections")
)

func main() {
	flag.Parse()
	if *flagClient {
		NewClient(*flagIn, *flagOut)
	} else if *flagServer {
		NewServer(*flagIn, *flagOut)
	} else {
		usage()
	}
	<-(chan int)(nil)
}
