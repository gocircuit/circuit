// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"github.com/gocircuit/circuit/use/n"
	"github.com/pkg/errors"
	// "bytes"
	"io"
	"os"

	"github.com/gocircuit/circuit/client"

	"github.com/urfave/cli"
)

func stack(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("recv needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	switch u := c.Walk(w).Get().(type) {
	case client.Server:
		r, err := u.Profile("goroutine")
		if err != nil {
			return errors.Wrapf(err, "error: %v", err)
		}
		io.Copy(os.Stdout, r)
	default:
		return errors.New("not a server")
	}
	return
}

func suicide(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("suicide needs one server anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(client.Server)
	if !ok {
		return errors.New("not a server")
	}
	u.Suicide()
	return
}

func join(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 2 {
		return errors.New("join needs one anchor argument and one circuit address argument")
	}
	// Verify the target circuit address is valid
	if _, err = n.ParseAddr(args[1]); err != nil {
		return errors.Wrapf(err, "argument %q is not a valid circuit address", args[1])
	}
	//
	w, _ := parseGlob(args[0])
	switch u := c.Walk(w).Get().(type) {
	case client.Server:
		if err = u.Rejoin(args[1]); err != nil {
			return errors.Wrapf(err, "error: %v", err)
		}
	default:
		return errors.New("not a server")
	}
	return
}
