// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"flag"
	"strconv"

	"github.com/gocircuit/circuit/client"
)

var flagMount = flag.String("m", "/circuit", "circuit mount point")

func main() {
	flag.Parse()
	c := client.Attach(*flagMount)
	ids := c.Circuits()
	waitexit := make(chan int, len(ids))
	for i, id := range ids {
		q := c.Circuit(id)
		a := q.Term("tutorial", "waitall")
		p := a.Proc("sleeper")
		p.Start(client.Command{
			Path: "sleep",
			Args: []string{strconv.Itoa(3*i)},
		})
		println("ping")
		go func() {
			p.Wait()
			waitexit <- 1
		}()
	}
	for _ = range ids {
		<-waitexit
		println("pong")
	}
}
