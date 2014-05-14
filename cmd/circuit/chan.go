// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	//"log"
	"io"
	"os"
	"strconv"

	"github.com/gocircuit/circuit/client"

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
	a := c.Walk(w)
	n, err := strconv.Atoi(args[1])
	if err != nil || n < 0 {
		fatalf("second argument to mkchan must be a non-negative integral capacity")
	}
	if _, err = a.MakeChan(n); err != nil {
		fatalf("mkchan error: %s", err)
	}
}

func send(x *cli.Context) {
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("send needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(client.Chan)
	if !ok {
		fatalf("not a channel")
	}
	msgw, err := u.Send()
	if err != nil {
		fatalf("send error: %v", err)
	}
	if _, err = io.Copy(msgw, os.Stdin); err != nil {
		fatalf("transmission error: %v", err)
	}
}

func recv(x *cli.Context) {
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("recv needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(client.Chan)
	if !ok {
		fatalf("not a channel")
	}
	msgr, err := u.Recv()
	if err != nil {
		fatalf("recv error: %v", err)
	}
	io.Copy(os.Stdout, msgr)
	// if _, err = io.Copy(os.Stdout, msgr); err != nil {
	// 	fatalf("transmission error: %v", err)
	// }
}

func clos(x *cli.Context) {
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("close needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(client.Chan)
	if !ok {
		fatalf("not a channel")
	}
	if err := u.Close(); err != nil {
		fatalf("close error: %v", err)
	}
}
