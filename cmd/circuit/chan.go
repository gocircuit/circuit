// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"log"
	"strconv"

	"github.com/codegangsta/cli"
)

// circuit mkchan /X1234/hola/charlie 0
func mkchan(x *cli.Context) {
	c := dial(x)
	args := x.Args()
	if len(args) != 2 {
		fatalf("mkchan needs an anchor and a capacity arguments")
	}
	w, _ := parseGlob(args[0])
	log.Printf("walk %v", w)
	a := c.Walk(w)
	n, err := strconv.Atoi(args[1])
	if err != nil || n < 0 {
		fatalf("second argument to mkchan must be a non-negative integral capacity")
	}
	if _, err = a.MakeChan(n); err != nil {
		fatalf("mkchan error: %s", err)
	}
}

func send(x *cli.Context) {}
func recv(x *cli.Context) {}
func clos(x *cli.Context) {}
func scrb(x *cli.Context) {}