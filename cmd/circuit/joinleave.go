// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"github.com/urfave/cli"
)

// circuit mk@join /X1234/hola/listy
func mkonjoin(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("mk@join needs an anchor argument")
	}
	w, _ := parseGlob(args[0])
	_, err := c.Walk(w).MakeOnJoin()
	if err != nil {
		fatalf("mk@join error: %s", err)
	}
}

// circuit mk@leave /X1234/hola/listy
func mkonleave(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("mk@leave needs an anchor argument")
	}
	w, _ := parseGlob(args[0])
	_, err := c.Walk(w).MakeOnLeave()
	if err != nil {
		fatalf("mk@leave error: %s", err)
	}
}
