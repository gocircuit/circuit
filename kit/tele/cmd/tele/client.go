// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"fmt"
	"net"
	"os"

	"github.com/gocircuit/circuit/kit/tele"
	"github.com/gocircuit/circuit/kit/tele/blend"
	"github.com/gocircuit/circuit/kit/tele/tcp"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

// Client
type Client struct {
	frame   trace.Frame
	tele    *blend.Transport
	inAddr  string
	outAddr string
}

func NewClient(inAddr, outAddr string) {
	cli := &Client{frame: trace.NewFrame("tele", "client"), outAddr: outAddr}
	cli.frame.Bind(cli)

	// Make teleport transport
	t := tele.NewStructOverTCP()

	// Listen on input TCP address
	l, err := net.Listen("tcp", inAddr)
	if err != nil {
		cli.frame.Printf("listen on teleport address %s (%s)", inAddr, err)
		os.Exit(1)
	}
	if inAddr == "" {
		inAddr = l.Addr().String()
		fmt.Println(inAddr)
	}
	cli.tele, cli.inAddr = t, inAddr
	go cli.loop(l)
	return
}

func (cli *Client) loop(l net.Listener) {
	for {
		inConn, err := l.Accept()
		if err != nil {
			cli.frame.Printf("accept on tcp address %s (%s)", cli.inAddr, err)
			os.Exit(1)
		}
		// Contact teleport server
		ds, _ := cli.tele.DialSession(tcp.Addr(cli.outAddr))
		tele := ds.Dial()
		// Write an empty chunk to mark the beginning of connection
		if err = tele.Write(&cargo{}); err != nil {
			cli.frame.Printf("first write (%s)", err)
			tele.Close()
			inConn.Close()
			continue
		}
		// Begin proxying
		Proxy(inConn, tele)
	}
}
