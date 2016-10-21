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
	"os"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/docker"
	"github.com/pkg/errors"

	"github.com/urfave/cli"
)

// circuit mkproc /X1234/hola/charlie << EOF
// { … }
// EOF
// TODO: Proc element disappears if command misspelled and error condition not obvious.
func mkproc(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("mkproc needs an anchor argument")
	}
	w, _ := parseGlob(args[0])
	buf, _ := ioutil.ReadAll(os.Stdin)
	var cmd client.Cmd
	if err = json.Unmarshal(buf, &cmd); err != nil {
		return errors.Wrapf(err, "command json not parsing: %v", err)
	}
	if x.Bool("scrub") {
		cmd.Scrub = true
	}
	p, err := c.Walk(w).MakeProc(cmd)
	if err != nil {
		return errors.Wrapf(err, "mkproc error: %s", err)
	}
	ps := p.Peek()
	if ps.Exit != nil {
		return errors.Errorf("%v", ps.Exit)
	}
	return
}

func mkdkr(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("mkdkr needs an anchor argument")
	}
	w, _ := parseGlob(args[0])
	buf, _ := ioutil.ReadAll(os.Stdin)
	var run docker.Run
	if err = json.Unmarshal(buf, &run); err != nil {
		return errors.Wrapf(err, "command json not parsing: %v", err)
	}
	if x.Bool("scrub") {
		run.Scrub = true
	}
	if _, err = c.Walk(w).MakeDocker(run); err != nil {
		return errors.Wrapf(err, "mkdkr error: %s", err)
	}
	return
}

// circuit signal kill /X1234/hola/charlie
func sgnl(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 2 {
		return errors.New("signal needs an anchor and a signal name arguments")
	}
	w, _ := parseGlob(args[1])
	u, ok := c.Walk(w).Get().(interface {
		Signal(string) error
	})
	if !ok {
		return errors.New("anchor is not a process or a docker container")
	}
	if err = u.Signal(args[0]); err != nil {
		return errors.Wrapf(err, "signal error: %v", err)
	}
	return
}

func wait(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("wait needs one anchor argument")
	}
	w, _ := parseGlob(args[0])

	var stat interface{}
	switch u := c.Walk(w).Get().(type) {
	case client.Proc:
		stat, err = u.Wait()
	case docker.Container:
		stat, err = u.Wait()
	default:
		return errors.New("anchor is not a process or a docker container")
	}
	if err != nil {
		return errors.Wrapf(err, "wait error: %v", err)
	}
	buf, _ := json.MarshalIndent(stat, "", "\t")
	fmt.Println(string(buf))
	return
}
