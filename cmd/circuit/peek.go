// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"fmt"
	"encoding/json"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/docker"

	"github.com/gocircuit/circuit/github.com/codegangsta/cli"
)

// circuit peek /X1234/hola/charlie
func peek(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("peek needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	switch t := c.Walk(w).Get().(type) {
	case client.Server:
		buf, _ := json.MarshalIndent(t.Peek(), "", "\t")
		fmt.Println(string(buf))
	case client.Chan:
		buf, _ := json.MarshalIndent(t.Stat(), "", "\t")
		fmt.Println(string(buf))
	case client.Proc:
		buf, _ := json.MarshalIndent(t.Peek(), "", "\t")
		fmt.Println(string(buf))
	case docker.Container:
		stat, err := t.Peek()
		if err != nil {
			fatalf("%v", err)
		}
		buf, _ := json.MarshalIndent(stat, "", "\t")
		fmt.Println(string(buf))
	case client.Subscription:
		buf, _ := json.MarshalIndent(t.Peek(), "", "\t")
		fmt.Println(string(buf))
	case nil:
		buf, _ := json.MarshalIndent(nil, "", "\t")
		fmt.Println(string(buf))
	default:
		fatalf("unknown element")
	}
}

func scrb(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("scrub needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	c.Walk(w).Scrub()
}
