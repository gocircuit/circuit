// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/github.com/codegangsta/cli"
)

func fatalf(format string, arg ...interface{}) {
	println(fmt.Sprintf(format, arg...))
	os.Exit(1)
}

func dial(x *cli.Context) *client.Client {
	var dialAddr string
	switch {
	case x.IsSet("dial"):
		dialAddr = x.String("dial")
	case os.Getenv("CIRCUIT") != "":
		buf, err := ioutil.ReadFile(os.Getenv("CIRCUIT"))
		if err != nil {
			fatalf("circuit environment file %s is not readable: %v", os.Getenv("CIRCUIT"), err)
		}
		dialAddr = strings.TrimSpace(string(buf))
	default:
		buf, err := ioutil.ReadFile(".circuit")
		if err != nil {
			fatalf("no dial address available; use flag -d or set CIRCUIT to a file name")
		}
		dialAddr = strings.TrimSpace(string(buf))
	}
	defer func() {
		if r := recover(); r != nil {
			fatalf("addressed server is gone or a newer one is in place")
		}
	}()
	return client.Dial(dialAddr)
}

// circuit ls /Q123/apps/charlie
// circuit ls /...
func ls(x *cli.Context) {
	c := dial(x)
	args := x.Args()
	if len(args) != 1 {
		println("ls needs a glob argument")
		os.Exit(1)
	}
	w, ellipses := parseGlob(args[0])
	list(0, "/", c.Walk(w), ellipses)
}

func list(level int, prefix string, anchor client.Anchor, recurse bool) {
	if anchor == nil {
		return
	}
	//println(fmt.Sprintf("prefix=%v a=%v/%T r=%v", prefix, anchor, anchor, recurse))
	var c children
	for n, a := range anchor.View() {
		e := &entry{n: n, a: a}
		v := a.Get()
		switch v.(type) {
		case client.Chan:
			e.k = "chan"
		case client.Proc:
			e.k = "proc"
		default:
			if level == 0 {
				e.k = "----"
			}
		}
		c = append(c, e)
	}
	sort.Sort(c)
	for _, e := range c {
		if e.k != "" {
			fmt.Printf("%4s %s%s\n", e.k, prefix, e.n)
		}
		if recurse {
			list(level + 1, prefix + e.n + "/", e.a, true)
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
	if walk[len(walk) - 1] == "..." {
		walk = walk[:len(walk)-1]
		ellipses = true
	}
	return
}
