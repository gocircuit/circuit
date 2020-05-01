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
	"path"
	"path/filepath"
	"os"
	"time"

	"github.com/hoijui/circuit/pkg/client"
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
func waitFotPayloadDeath(c *client.Client, payloadAnchor string) (recov interface{}) {
	// defer func() { // catch panics caused by unexpected death of the server hosting the payload
	// 	recov = recover()
	// }()
	t := c.Walk(client.Split(payloadAnchor)) // Access the process anchor of the currently-running payload of the virus.
	t.Get().(client.Proc).Wait() // Wait until the payload process exits.
	t.Scrub() // scrub payload anchor from old process element
	time.Sleep(2*time.Second) // Wait a touch to slow down the spin
	return
}

// The initial invocation of the virus:
//	virus DIALIN_CIRCUIT
// To invoke the virus in the role of a nucleus process:
// 	virus DIALIN_CIRCUIT PAYLOAD_ANCHOR SELF_ANCHOR
//
func main() {
	var payloadAnchor, nucleusAnchor string
	switch len(os.Args) {
	case 2: // initial command-line invocation
	case 4: // invocation in role of nucleus
		payloadAnchor = os.Args[2]
		nucleusAnchor = os.Args[3]
	default:
		println("usage: virus circuit://...")
		os.Exit(1)
	}
	println("virus dialing into", os.Args[1])
	c := client.Dial(os.Args[1], nil)

	// The nucleus role waits for the payload process to die before it proceeds.
	if nucleusAnchor != "" {
		waitFotPayloadDeath(c, payloadAnchor)
		c.Walk(client.Split(nucleusAnchor)).Scrub() // remove anchor pointing to us
	}
	spawnNucleus(c, spawnPayload(c))
}

func spawnPayload(c *client.Client) (payloadAnchor string) {
	service := client.Cmd{
		Path: "/usr/bin/say", // say is a standard OSX command which speaks, so it's easy to hear the virus in action.
		Args: []string{"i am a virus"},
	}
	a := pickServer(c) // Randomly choose a circuit server to host the virus payload.
	pservice, err := a.Walk([]string{"virus", "payload"}).MakeProc(service) // Run the payload process
	if err != nil {
		println("payload not created:", err.Error())
		os.Exit(1)
	}
	if err := pservice.Peek().Exit; err != nil {
		println("payload not started:", err.Error())
		os.Exit(1)
	}
	pservice.Stdin().Close() // Close the standard input of the payload to indicate no intention to write data.
	return path.Join("/", a.ServerID(), "virus", "payload")
}

func spawnNucleus(c *client.Client, payloadAnchor string) {
	b := pickServer(c)
	virus, _ := filepath.Abs(os.Args[0]) // We assume that the virus binary is on the same path everywhere
	nucleusAnchor := path.Join("/", b.ServerID(), "virus", "nucleus")
	nucleus := client.Cmd{
		Path: virus,
		Args: []string{
			b.Addr(), // dial-in circuit server address
			payloadAnchor, // payload anchor
			nucleusAnchor, // anchor of the spawned nucleus itself
		},
	}
	pnucleus, err := b.Walk([]string{"virus", "nucleus"}).MakeProc(nucleus)
	if err != nil {
		println("nucleus not created:", err.Error())
		os.Exit(1)
	}
	if err := pnucleus.Peek().Exit; err != nil {
		println("nucleus not started:", err.Error())
		os.Exit(1)
	}
	pnucleus.Stdin().Close()
}
