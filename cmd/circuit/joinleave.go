// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// circuit mk@join /X1234/hola/listy
func mkonjoin(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("mk@join needs an anchor argument")
	}
	w, _ := parseGlob(args[0])
	if _, err = c.Walk(w).MakeOnJoin(); err != nil {
		return errors.Wrapf(err, "mk@join error: %s", err)
	}
	return
}

// circuit mk@leave /X1234/hola/listy
func mkonleave(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("mk@leave needs an anchor argument")
	}
	w, _ := parseGlob(args[0])
	if _, err = c.Walk(w).MakeOnLeave(); err != nil {
		return errors.Wrapf(err, "mk@leave error: %s", err)
	}
	return
}
