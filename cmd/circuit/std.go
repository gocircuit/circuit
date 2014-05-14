// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"io"
	"os"
	"github.com/gocircuit/circuit/client"

	"github.com/codegangsta/cli"
)

func stdin(x *cli.Context) {
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("stdin needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(client.Proc)
	if !ok {
		fatalf("not a process")
	}
	q := u.Stdin()
	if _, err := io.Copy(q, os.Stdin); err != nil {
		fatalf("transmission error: %v", err)
	}
	if err := q.Close(); err != nil {
		fatalf("error closing stdin: %v", err)
	}
}

func stdout(x *cli.Context) {
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("stdout needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(client.Proc)
	if !ok {
		fatalf("not a process")
	}
	io.Copy(os.Stdout, u.Stdout())
	// if _, err := io.Copy(os.Stdout, u.Stdout()); err != nil {
	// 	fatalf("transmission error: %v", err)
	// }
}

func stderr(x *cli.Context) {
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("stderr needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(client.Proc)
	if !ok {
		fatalf("not a process")
	}
	io.Copy(os.Stdout, u.Stderr())
	// if _, err := io.Copy(os.Stdout, u.Stderr()); err != nil {
	// 	fatalf("transmission error: %v", err)
	// }
}
