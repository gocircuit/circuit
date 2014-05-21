// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// This package provides the executable program for the resource-sharing circuit app
package main

import (
	"os"
	"github.com/gocircuit/circuit/github.com/codegangsta/cli"
)

func main() {
	
	app := cli.NewApp()
	app.Name = "circuit"
	app.Usage = "Circuit server and client tool"
	app.Commands = []cli.Command{
		// circuit
		{
			Name: "start",
			Usage: "Run a circuit worker on this machine",
			Action: server,
			Flags: []cli.Flag{
				cli.StringFlag{"addr, a", "", "address of circuit server"},
				cli.StringFlag{"mutex, m", "", "directory to use as a circuit instance mutex lock"},
				cli.StringFlag{"join, j", "", "join a circuit through a current member by url"},
			},
	 	},
		{
			Name: "ls",
			Usage: "List circuit elements",
			Action: ls,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
				cli.BoolFlag{"long, l", "show detailed anchor information"},
				cli.BoolFlag{"depth, de", "traverse anchors in depth-first order (leaves first)"},
			},
		},
		// channel-specific
		{
			Name: "mkchan",
			Usage: "Create a channel element",
			Action: mkchan,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
		{
			Name: "send",
			Usage: "Send data to the channel from standard input",
			Action: send,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
		{
			Name: "recv",
			Usage: "Receive data from the channel on stadard output",
			Action: recv,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
		{
			Name: "close",
			Usage: "Close the channel after all current transmissions complete",
			Action: clos,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
		// common
		{
			Name: "scrub",
			Usage: "Abort and remove an element",
			Action: scrb,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
		{
			Name: "peek",
			Usage: "Query element state asynchronously",
			Action: peek,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
		// proc-specific
		{
			Name: "mkproc",
			Usage: "Create a process element",
			Action: mkproc,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
		{
			Name: "signal",
			Usage: "Send a signal to a running process",
			Action: sgnl,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
		{
			Name: "wait",
			Usage: "Wait until a process exits",
			Action: wait,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
		// stdin, stdout, stderr
		{
			Name: "stdin",
			Usage: "Forward this tool's standard input to that of the process",
			Action: stdin,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
		{
			Name: "stdout",
			Usage: "Forward the standard output of the process to the standard output of this tool",
			Action: stdout,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
		{
			Name: "stderr",
			Usage: "Forward the standard error of the process to the standard output of this tool",
			Action: stderr,
			Flags: []cli.Flag{
				cli.StringFlag{"dial, d", "", "circuit member to dial into"},
			},
		},
	}
	app.Run(os.Args)
}
