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

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func stdin(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("stdin needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(interface {
		Stdin() io.WriteCloser
	})
	if !ok {
		return errors.New("not a process or a container")
	}
	q := u.Stdin()
	if _, err = io.Copy(q, os.Stdin); err != nil {
		return errors.Wrapf(err, "transmission error: %v", err)
	}
	if err = q.Close(); err != nil {
		return errors.Wrapf(err, "error closing stdin: %v", err)
	}
	return
}

func stdout(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("stdout needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(interface {
		Stdout() io.ReadCloser
	})
	if !ok {
		return errors.New("not a process or a container")
	}
	io.Copy(os.Stdout, u.Stdout())
	return
}

func stderr(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		return errors.New("stderr needs one anchor argument")
	}
	w, _ := parseGlob(args[0])
	u, ok := c.Walk(w).Get().(interface {
		Stderr() io.ReadCloser
	})
	if !ok {
		return errors.New("not a process or a container")
	}
	io.Copy(os.Stdout, u.Stderr())
	// if _, err := io.Copy(os.Stdout, u.Stderr()); err != nil {
	// 	fatalf("transmission error: %v", err)
	// }
	return
}
