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

	"github.com/codegangsta/cli"
)

// circuit peek /X1234/hola/charlie
func peek(x *cli.Context) {
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("stat needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	switch t := c.Walk(w).Get().(type) {
	case client.Chan:
		buf, _ := json.MarshalIndent(t.Stat(), "", "\t")
		fmt.Println(string(buf))
	case client.Proc:
		buf, _ := json.MarshalIndent(t.Peek(), "", "\t")
		fmt.Println(string(buf))
	case nil:
		buf, _ := json.MarshalIndent(nil, "", "\t")
		fmt.Println(string(buf))
	default:
		fatalf("unknown element")
	}
}
