// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gocircuit/circuit/client"
	"github.com/pkg/errors"

	"github.com/urfave/cli"
)

func recv(x *cli.Context) (err error) {
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
	case client.Chan:
		msgr, err := u.Recv()
		if err != nil {
			return errors.Wrapf(err, "recv error: %v", err)
		}
		io.Copy(os.Stdout, msgr)
	case client.Subscription:
		v, ok := u.Consume()
		if !ok {
			return errors.New("eof")
		}
		fmt.Println(v)
		os.Stdout.Sync()
	default:
		return errors.New("not a channel or subscription")
	}
	return
}
