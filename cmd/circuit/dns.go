// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"github.com/gocircuit/circuit/client"
	"github.com/pkg/errors"

	"github.com/urfave/cli"
)

func mkdns(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) < 1 {
		return errors.New("mkdns needs an anchor and an optional address arguments")
	}
	var addr string
	if len(args) == 2 {
		addr = args[1]
	}
	w, _ := parseGlob(args[0])

	if _, err = c.Walk(w).MakeNameserver(addr); err != nil {
		return errors.Wrapf(err, "mkdns error: %s", err)
	}
	return
}

func nset(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 2 {
		return errors.New("set needs an anchor and a resource record arguments")
	}
	w, _ := parseGlob(args[0])
	switch u := c.Walk(w).Get().(type) {
	case client.Nameserver:
		err := u.Set(args[1])
		if err != nil {
			return errors.Wrapf(err, "set resoure record error: %v", err)
		}
	default:
		return errors.New("not a nameserver element")
	}
	return
}

func nunset(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 2 {
		return errors.New("unset needs an anchor and a resource name arguments")
	}
	w, _ := parseGlob(args[0])
	switch u := c.Walk(w).Get().(type) {
	case client.Nameserver:
		u.Unset(args[1])
	default:
		return errors.New("not a nameserver element")
	}
	return
}
