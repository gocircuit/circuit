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
	"github.com/pkg/errors"

	"github.com/urfave/cli"
)

// circuit mkchan /X1234/hola/charlie 0
func mkchan(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 2 {
		return errors.New("mkchan needs an anchor and a capacity arguments")
	}
	w, _ := parseGlob(args[0])
	a := c.Walk(w)
	n, err := strconv.Atoi(args[1])
	if err != nil || n < 0 {
		return errors.New("second argument to mkchan must be a non-negative integral capacity")
	}
	if _, err = a.MakeChan(n); err != nil {
		return errors.Wrapf(err, "mkchan error: %s", err)
	}
	return
}

func send(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("send needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(client.Chan)
	if !ok {
		return errors.New("not a channel")
	}
	msgw, err := u.Send()
	if err != nil {
		return errors.Wrapf(err, "send error: %v", err)
	}
	if _, err = io.Copy(msgw, os.Stdin); err != nil {
		return errors.Wrapf(err, "transmission error: %v", err)
	}
	return
}

func clos(x *cli.Context) (err error) {
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("close needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(client.Chan)
	if !ok {
		return errors.New("not a channel")
	}
	if err := u.Close(); err != nil {
		return errors.Wrapf(err, "close error: %v", err)
	}
	return
}
