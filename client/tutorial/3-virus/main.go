// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

/*

	Virus is a simple, mildly-resilient to failure mechanism that runs around a cluster
	and delivers its payload (a process execution), in a self-sustained fashion.

	The virus mechanism contains two parts: a payload process and a nucleus process.

	The payload can be any OS process, available locally as a binary.

	The nucleus is embodied in this circuit application. It executes the payload OS process
	on a random machine in the circuit cluster. Then it installs itself on a machine different
	from that of the payload process, and proceeds to watch the payload until it dies.

	When the payload dies, the nucleus executes a new payload instance on another 
	randomly chosen host, and replaces itself with a new nucleus process on yet
	another new random host. And so on ...

*/
package main

import (
	"path/filepath"
	"os"
	"time"

	"github.com/gocircuit/circuit/client"
)

// pickServer returns the root anchor of a randomly-chosen circuit server in the cluster.
func pickServer(c *client.Client) client.Anchor {
	for _, r := range c.View() {
		return r
	}
	panic(0)
}

// waitForPayloadDeath blocks until the payload process stored at anchor exits, for whatever reason.
// anchor is the anchor path used by the virus logic.
func waitFotPayloadDeath(c *client.Client, anchor string) {

	// catch panics caused by unexpected death of the server hosting the payload
	defer func() {
		recover()
	}()

	// Access the process anchor that started this very process and
	// remove it to make room for the new one.
	// Note that scrubbing a process anchor removes the process element,
	// but in no way affects the underlying OS process.
	walkToVirus := client.Split(anchor)
	c.Walk(append(walkToVirus, "nucleus")).Scrub()

	// Access the process anchor of the currently-running payload of the virus.
	t := c.Walk(append(walkToVirus, "payload"))
	// Wait until the payload process exits.
	t.Get().(client.Proc).Wait()
	// Remove the anchor of the now-dead payload process.
	t.Scrub()

	// Wait a touch to prevent spinning, if the payload exits immediately every time it is run.
	time.Sleep(time.Second/2)
}

func main() {

	// The first argument is the circuit server address that this execution will use.
	c := client.Dial(os.Args[1])

	// 
	if len(os.Args) == 3 {
		waitFotPayloadDeath(c, os.Args[2])
	}

	// Start the payload process
	service := client.Cmd{
		Path: "/usr/bin/say", // say is a standard OSX command which speaks, so it's easy to hear the virus in action.
		Args: []string{"i am a virus"},
	}
	// Randomly choose a circuit server to host the virus payload.
	a := pickServer(c)
	// Run the payload
	pservice, _ := a.Walk([]string{"virus", "payload"}).MakeProc(service)
	if err := pservice.Peek().Exit; err != nil {
		println("payload not started:", err.Error())
		return
	}
	// Close the standard input of the virus to indicate no intention to write data.
	pservice.Stdin().Close()

	// Start the virus nucleus process, which will wait until the payload completes,
	// and then start a payload as well as a new nucleus elsewhere, over and over again.
	b := pickServer(c)
	virus, _ := filepath.Abs(os.Args[0]) // We assume that the virus binary is on the same path everywhere
	nucleus := client.Cmd{
		Path: virus,
		Args: []string{b.Addr(), "/" + a.ServerID() + "/virus"},
	}
	pnucleus, _ := b.Walk([]string{"virus", "nucleus"}).MakeProc(nucleus)
	if err := pnucleus.Peek().Exit; err != nil {
		println("nucleus not started:", err.Error())
		return
	}
	pnucleus.Stdin().Close()
}
