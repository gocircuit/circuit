// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocircuit/circuit/client"
	"github.com/codegangsta/cli"
)

func fatalf(format string, arg ...interface{}) {
	fmt.Fprintf(os.Stderr, format, arg...)
	os.Exit(1)
}

func dial(x *cli.Context) *client.Client {
	var dialAddr string
	switch {
	case x.IsSet("dial"):
		dialAddr = x.String("dial")
	case os.Getenv("CIRCUIT_DIAL") != "":
		dialAddr = os.Getenv("CIRCUIT_DIAL")
	default:
		panic("no dialAddr address available")
	}
	return client.Dial(dialAddr)
}

// circuit ls /Q123/apps/charlie
// circuit ls /...
func ls(x *cli.Context) {
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		println("ls needs a glob argument")
		os.Exit(1)
	}
	w, ellipses := parseGlob(args[0])
	list("/", c.Walk(w), ellipses)
}

func list(prefix string, anchor client.Anchor, recurse bool) {
	for n, a := range anchor.View() {
		v := a.Get()
		var k string
		switch v.(type) {
		case client.Chan:
			k = "chan"
		case client.Proc:
			k = "proc"
		default:
			k = "····"
		}
		fmt.Printf("%4s %s%s\n", k, prefix, n)
		if recurse {
			list(prefix + n + "/", a, true)
		}
	}
}

func parseGlob(pattern string) (walk []string, ellipses bool) {
	for _, p := range strings.Split(pattern, "/") {
		if len(p) == 0 {
			continue
		}
		walk = append(walk, p)
	}
	if len(walk) == 0 {
		return
	}
	if walk[len(walk) - 1] == "..." {
		walk = walk[:len(walk)-1]
		ellipses = true
	}
	return
}
