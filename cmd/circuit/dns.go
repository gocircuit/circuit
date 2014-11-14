// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"github.com/gocircuit/circuit/client"

	"github.com/gocircuit/circuit/github.com/codegangsta/cli"
)

func mkdns(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) < 1 {
		fatalf("mkdns needs an anchor and an optional address arguments")
	}
	var addr string
	if len(args) == 2 {
		addr = args[1]
	}
	w, _ := parseGlob(args[0])
	_, err := c.Walk(w).MakeNameserver(addr)
	if err != nil {
		fatalf("mkdns error: %s", err)
	}
}

func nset(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 2 {
		fatalf("set needs an anchor and a resource record arguments")
	}
	w, _ := parseGlob(args[0])
	switch u := c.Walk(w).Get().(type) {
	case client.Nameserver:
		err := u.Set(args[1])
		if err != nil {
			fatalf("set resoure record error: %v", err)
		}
	default:
		fatalf("not a nameserver element")
	}
}

func nunset(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 2 {
		fatalf("unset needs an anchor and a resource name arguments")
	}
	w, _ := parseGlob(args[0])
	switch u := c.Walk(w).Get().(type) {
	case client.Nameserver:
		u.Unset(args[1])
	default:
		fatalf("not a nameserver element")
	}
}
