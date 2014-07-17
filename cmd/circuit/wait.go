// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"bufio"
	"os"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/docker"
	"github.com/gocircuit/circuit/github.com/codegangsta/cli"
)

func waitall(x *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fatalf("error: %v", r)
		}
	}()
	c := dial(x)
	s := bufio.NewScanner(os.Stdin)
	var t []string
	ch := make(chan int)
	for s.Scan() {
		src := s.Text()
		i := len(t)
		t = append(t, src)
		go func() { // wait on w
			w, _ := parseGlob(src)
			var err error
			switch u := c.Walk(w).Get().(type) {
			case client.Proc:
				_, err = u.Wait()
			case docker.Container:
				_, err = u.Wait()
			default:
				println("anchor", w, " is not a process or a docker container")
			}
			if err != nil {
				fatalf("wait error: %v", err)
			}
			ch <- i
		}()
	}
	for _ = range t {
		<-ch
	}
}
