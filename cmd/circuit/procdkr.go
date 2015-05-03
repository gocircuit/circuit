// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/mcqueenorama/circuit/client"
	"github.com/mcqueenorama/circuit/client/docker"

	"github.com/mcqueenorama/circuit/github.com/codegangsta/cli"
)

//make this timeout come from the json input payload
var timeout = time.Duration(33)

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

	if len(args) != 1 && !x.Bool("all") {
		fatalf("runproc needs an anchor argument or use the --all flag to do the entire circuit")
	}
	buf, _ := ioutil.ReadAll(os.Stdin)
	var cmd client.Cmd
	if err := json.Unmarshal(buf, &cmd); err != nil {
		fatalf("command json not parsing: %v", err)
	}
	if x.Bool("scrub") {
		cmd.Scrub = true
	}

	el := string("")
	if len(cmd.Name) > 0 {
		el = cmd.Name
	} else {
		el = cmd.Path
	}

	if x.Bool("all") {

		w, _ := parseGlob("/")

		anchor := c.Walk(w)

		for _, a := range anchor.View() {
			w2, _ := parseGlob(a.Path() + "/" + el)
			a2 := c.Walk(w2)
			_runproc(a2, cmd, x.Bool("tag"))
		}

	} else {

		w, _ := parseGlob(args[0] + "/" + el)
		a := c.Walk(w)
		_runproc(a, cmd, x.Bool("tag"))

	}

}

func _runproc(a client.Anchor, cmd client.Cmd, tags bool) {

	p, err := a.MakeProc(cmd)
	if err != nil {
		fatalf("mkproc error: %s", err)
	}

	q := p.Stdin()
	if err := q.Close(); err != nil {
		fatalf("error closing stdin: %v", err)
	}

	if tags {

		done := make(chan bool)

		go func(a string, r io.Reader, w io.Writer, d chan bool) {

			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				fmt.Fprintf(w, "%s %s\n", a, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				d <- false
			}

			d <- true

		}(a.Path(), p.Stdout(), os.Stdout, done)

		select {
		case rv := <-done:
			if !rv {
				fmt.Fprintln(os.Stderr, "error prefixing the data")
			}
		case <-time.After(time.Second * timeout):
			fmt.Fprintln(os.Stderr, "timeout waiting for data")
		}

	} else {

		io.Copy(os.Stdout, p.Stdout())

	}

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
	u, ok := c.Walk(w).Get().(interface {
		Signal(string) error
	})
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
