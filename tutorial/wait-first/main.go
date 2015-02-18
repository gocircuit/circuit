// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

// WaitFirst is a simple circuit application that starts a few worker processes and waits until someone finishes.
package main

import (
	"os"
	"strconv"

	"github.com/gocircuit/circuit/client"
)

const n = 5

// pickServer returns the root anchor of a randomly-chosen circuit server in the cluster.
func pickServer(c *client.Client) client.Anchor {
	for _, r := range c.View() {
		return r
	}
	panic(0)
}

func main() {

	// The first argument is the circuit server address that this execution will use.
	c := client.Dial(os.Args[1], nil)

	// Fire-off a few payload processes.
	ch := make(chan int)
	for i := 0; i < n; i++ {
		cmd := client.Cmd{
			Path: "/bin/sleep",
			Args: []string{strconv.Itoa(3+i*3)},
		}
		i_ := i
		go func() {
			// Pick a random circuit server to run payload on.
			t := pickServer(c).Walk([]string{"wait-first", strconv.Itoa(i_)})
			// Execute the process and store it in the anchor.
			p, _ := t.MakeProc(cmd)
			// Close the process standard input to indicte no intention to write data.
			p.Stdin().Close()
			// Block until the process exits.
			p.Wait()
			// Notify the higher-level synchronization logic
			ch <- 1
			// Remove the anchor storing the process.
			t.Scrub()
			println("Payload", i_+1, "finished.")
		}()
	}
	<-ch
	println("One done.")
}
