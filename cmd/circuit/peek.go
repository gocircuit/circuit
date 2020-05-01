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

	"github.com/hoijui/circuit/pkg/client"
	"github.com/hoijui/circuit/pkg/client/docker"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
)

// circuit peek /X1234/hola/charlie
func peek(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if x.NArg() != 1 {
		return errors.New("peek needs one anchor argument")
	}
	w, _ := parseGlob(args.Get(0))
	switch t := c.Walk(w).Get().(type) {
	case client.Server:
		buf, _ := json.MarshalIndent(t.Peek(), "", "\t")
		fmt.Println(string(buf))
	case client.Chan:
		buf, _ := json.MarshalIndent(t.Stat(), "", "\t")
		fmt.Println(string(buf))
	case client.Proc:
		buf, _ := json.MarshalIndent(t.Peek(), "", "\t")
		fmt.Println(string(buf))
	case client.Nameserver:
		buf, _ := json.MarshalIndent(t.Peek(), "", "\t")
		fmt.Println(string(buf))
	case docker.Container:
		stat, err := t.Peek()
		if err != nil {
			return errors.Wrapf(err, "%v", err)
		}
		buf, _ := json.MarshalIndent(stat, "", "\t")
		fmt.Println(string(buf))
	case client.Subscription:
		buf, _ := json.MarshalIndent(t.Peek(), "", "\t")
		fmt.Println(string(buf))
	case nil:
		buf, _ := json.MarshalIndent(nil, "", "\t")
		fmt.Println(string(buf))
	default:
		return errors.New("unknown element")
	}
	return
}

func scrb(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if x.NArg() != 1 {
		return errors.New("scrub needs one anchor argument")
	}
	w, _ := parseGlob(args.Get(0))
	c.Walk(w).Scrub()
	return
}
