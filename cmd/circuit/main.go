// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// This package provides the executable program for the resource-sharing circuit app
package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "circuit"
	app.Usage = "Circuit server and client tool"
	app.Commands = []*cli.Command{
		// circuit
		{
			Name:   "start",
			Usage:  "Run a circuit worker on this machine",
			Action: server,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "addr, a", Value: "0.0.0.0:0", Usage: "Address of circuit server."},
				&cli.StringFlag{Name: "if", Value: "", Usage: "Bind any available port on the specified interface."},
				&cli.StringFlag{Name: "var", Value: "", Usage: "Lock and log directory for the circuit server."},
				&cli.StringFlag{Name: "join, j", Value: "", Usage: "Join a circuit through a current member by address."},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File with HMAC credentials for HMAC/RC4 transport security.", EnvVars: []string{"CIRCUIT_HMAC"}},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.BoolFlag{Name: "docker", Usage: "Enable docker elements; docker command must be executable"},
			},
		},
		{
			Name:   "keygen",
			Usage:  "Generate a new random HMAC key",
			Action: keygen,
		},
		{
			Name:   "ls",
			Usage:  "List circuit elements",
			Action: ls,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.BoolFlag{Name: "long, l", Usage: "show detailed anchor information"},
				&cli.BoolFlag{Name: "depth, de", Usage: "traverse anchors in depth-first order (leaves first)"},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		// subscription-specific
		{
			Name:   "mk@join",
			Usage:  "Create a subscription element, receiving server join events",
			Action: mkonjoin,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "mk@leave",
			Usage:  "Create a subscription element, receiving server leave events",
			Action: mkonleave,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		// server-specific
		{
			Name:   "stk",
			Usage:  "Print the runtime stack trace of a server element",
			Action: stack,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "join",
			Usage:  "Merge the networks of this circuit server and that of the argument circuit address",
			Action: join,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "suicide",
			Usage:  "Kill a chosen circuit daemon",
			Action: stack,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		// channel-specific
		{
			Name:   "mkchan",
			Usage:  "Create a channel element",
			Action: mkchan,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "send",
			Usage:  "Send data to the channel from standard input",
			Action: send,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "recv",
			Usage:  "Receive data from a channel or a subscription on stadard output",
			Action: recv,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "close",
			Usage:  "Close the channel after all current transmissions complete",
			Action: clos,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		// common
		{
			Name:   "scrub",
			Usage:  "Abort and remove an element",
			Action: scrb,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "peek",
			Usage:  "Query element state asynchronously",
			Action: peek,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		// nameserver
		{
			Name:   "mkdns",
			Usage:  "Create a nameserver element",
			Action: mkdns,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "set",
			Usage:  "Set a resource record in a nameserver element",
			Action: nset,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "unset",
			Usage:  "Remove all resource records for a name in a nameserver element",
			Action: nunset,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		// proc/dkr-specific
		{
			Name:   "mkdkr",
			Usage:  "Create a docker container element",
			Action: mkdkr,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.BoolFlag{Name: "scrub", Usage: "scrub the process anchor automatically on exit"},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "mkproc",
			Usage:  "Create a process element",
			Action: mkproc,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.BoolFlag{Name: "scrub", Usage: "scrub the process anchor automatically on exit"},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "signal",
			Usage:  "Send a signal to a running process",
			Action: sgnl,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "wait",
			Usage:  "Wait until a process exits",
			Action: wait,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "waitall",
			Usage:  "Wait until a set of processes all exit",
			Action: waitall,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		// stdin, stdout, stderr
		{
			Name:   "stdin",
			Usage:  "Forward this tool's standard input to that of the process",
			Action: stdin,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "stdout",
			Usage:  "Forward the standard output of the process to the standard output of this tool",
			Action: stdout,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "stderr",
			Usage:  "Forward the standard error of the process to the standard output of this tool",
			Action: stderr,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial, d", Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
