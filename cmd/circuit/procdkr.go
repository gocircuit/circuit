// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"io"
	"os"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/docker"

	"github.com/gocircuit/circuit/github.com/codegangsta/cli"
)

// circuit mkproc /X1234/hola/charlie << EOF
// { … }
// EOF
// TODO: Proc element disappears if command misspelled and error condition not obvious.
func mkproc(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("mkproc needs an anchor argument")
	}
	w, _ := parseGlob(args[0])
	buf, _ := ioutil.ReadAll(os.Stdin)
	var cmd client.Cmd
	if err := json.Unmarshal(buf, &cmd); err != nil {
		fatalf("command json not parsing: %v", err)
	}
	if x.Bool("scrub") {
		cmd.Scrub = true
	}
	p, err := c.Walk(w).MakeProc(cmd)
	if err != nil {
		fatalf("mkproc error: %s", err)
	}
	ps := p.Peek()
	if ps.Exit != nil {
		fatalf("%v", ps.Exit)
	}
}

func runproc(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("runproc needs an anchor argument")
	}
	w, _ := parseGlob(args[0])
	buf, _ := ioutil.ReadAll(os.Stdin)
	var cmd client.Cmd
	if err := json.Unmarshal(buf, &cmd); err != nil {
		fatalf("command json not parsing: %v", err)
	}
	if x.Bool("scrub") {
		cmd.Scrub = true
	}
	a := c.Walk(w)
	p, err := a.MakeProc(cmd)
	if err != nil {
		fatalf("mkproc error: %s", err)
	}
	// ps := p.Peek()
	// if ps.Exit != nil {
	// 	fatalf("%v", ps.Exit)
	// }
	q := p.Stdin()
	if err := q.Close(); err != nil {
		fatalf("error closing stdin: %v", err)
	}
	io.Copy(os.Stdout, p.Stdout())
	a.Scrub()
}

func mkdkr(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("mkdkr needs an anchor argument")
	}
	w, _ := parseGlob(args[0])
	buf, _ := ioutil.ReadAll(os.Stdin)
	var run docker.Run
	if err := json.Unmarshal(buf, &run); err != nil {
		fatalf("command json not parsing: %v", err)
	}
	if x.Bool("scrub") {
		run.Scrub = true
	}
	_, err := c.Walk(w).MakeDocker(run)
	if err != nil {
		fatalf("mkdkr error: %s", err)
	}
}

// circuit signal kill /X1234/hola/charlie
func sgnl(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 2 {
		fatalf("signal needs an anchor and a signal name arguments")
	}
	w, _ := parseGlob(args[1])
	u, ok := c.Walk(w).Get().(interface{Signal(string) error})
	if !ok {
		fatalf("anchor is not a process or a docker container")
	}
	if err := u.Signal(args[0]); err != nil {
		fatalf("signal error: %v", err)
	}
}

func wait(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		fatalf("wait needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	//
	var stat interface{}
	var err error
	switch u := c.Walk(w).Get().(type) {
	case client.Proc:
		stat, err = u.Wait()
	case docker.Container:
		stat, err = u.Wait()
	default:
		fatalf("anchor is not a process or a docker container")
	}
	if err != nil {
		fatalf("wait error: %v", err)
	}
	buf, _ := json.MarshalIndent(stat, "", "\t")
	fmt.Println(string(buf))
}
