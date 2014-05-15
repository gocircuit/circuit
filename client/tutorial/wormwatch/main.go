// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"os"
	"strconv"

	"github.com/gocircuit/circuit/client"
)

// pick returns the root anchor of a randomly-chosen circuit member
func pick(c *client.Client) client.Anchor {
	for _, r := range c.View() {
		return r
	}
	panic(0)
}

// watch
func watch(c *client.Client, service string) {
	c.Walk(client.Split(service))
	??
}

// wormwatch dial_url service_anchor?
func main() {
	c := client.Dial(os.Args[1]) // argument is the url of a circuit server
	ch := make(chan int)

	service := client.Cmd{ // a pretend long-running user binary
		Path: "/bin/sleep",
		Args: []string{strconv.Itoa(5)}, // with simulated unexpected exits
	}

	a := pick(c)
	watch := client.Cmd{
		Path: os.Args[0], // we assume that the binary of this tool is on the same path everywhere
		Args: []string{c.Addr()}, // instruct the watcher to use its local circuit for dial-in
	}

	t := a.Walk([]string{"worm_watch", "service"})
	pservice, _ := t.MakeProc(service)
	pservice.Stdin().Close()
	p.Wait()

	ch <- 1
	println("process", i_+1, "done")
}
