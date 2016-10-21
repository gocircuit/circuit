// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/docker"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// circuit ls /Q123/apps/charlie
// circuit ls /...
func ls(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		println("ls needs a glob argument")
		os.Exit(1)
	}
	w, ellipses := parseGlob(args[0])
	list(0, "/", c.Walk(w), ellipses, x.Bool("long"), x.Bool("depth"))
	return
}

func list(level int, prefix string, anchor client.Anchor, recurse, long, depth bool) {
	if anchor == nil {
		return
	}
	// println(fmt.Sprintf("prefix=%v a=%v/%T r=%v", prefix, anchor, anchor, recurse))
	var c children
	for n, a := range anchor.View() {
		e := &entry{n: n, a: a}
		v := a.Get()
		switch t := v.(type) {
		case client.Server:
			e.k = "server"
		case client.Chan:
			e.k = "chan"
		case client.Proc:
			e.k = "proc"
			// if t.GetCmd().Scrub {
			// 	e.k = "proc-autoscrub"
			// } else {
			// 	e.k = "proc"
			// }
		case client.Nameserver:
			e.k = "dns"
		case docker.Container:
			e.k = "docker"
		case client.Subscription:
			e.k = "@" + t.Peek().Source
		default:
			e.k = "·"
		}
		c = append(c, e)
	}
	sort.Sort(c)
	for _, e := range c {
		if recurse && depth {
			list(level+1, prefix+e.n+"/", e.a, true, long, depth)
		}
		if long {
			fmt.Printf("%-15s %s%s\n", e.k, prefix, e.n)
		} else {
			fmt.Printf("%s%s\n", prefix, e.n)
		}
		if recurse && !depth {
			list(level+1, prefix+e.n+"/", e.a, true, long, depth)
		}
	}
}

type entry struct {
	n string
	a client.Anchor
	k string
}

type children []*entry

func (c children) Len() int {
	return len(c)
}

func (c children) Less(i, j int) bool {
	return c[i].n < c[j].n
}

func (c children) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func parseGlob(pattern string) (walk []string, ellipses bool) {
	for _, p := range strings.Split(pattern, "/") {
		if len(p) == 0 {
			continue
		}
		walk = append(walk, p)
	}
	if len(walk) == 0 {
		return
	}
	if walk[len(walk)-1] == "..." {
		walk = walk[:len(walk)-1]
		ellipses = true
	}
	return
}
