// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"bufio"
	"log"
	"os"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/docker"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func waitall(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error: %v", r)
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
			var e error
			switch u := c.Walk(w).Get().(type) {
			case client.Proc:
				_, e = u.Wait()
			case docker.Container:
				_, e = u.Wait()
			default:
				println("anchor", w, " is not a process or a docker container")
			}
			if e != nil {
				log.Fatal(errors.Errorf("wait error: %v", e))
			}
			ch <- i
		}()
	}
	for _ = range t {
		<-ch
	}
	return
}
