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
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "circuit"
	app.Usage = "Circuit server and client tool"
	app.Commands = []cli.Command{
		{
			Name: "server",
			ShortName: "srv",
			//Usage: "Run a circuit server",
			Action: server,
			Flags: []cli.Flag{
				cli.StringFlag{"addr, a", "", "address of circuit server"},
				cli.StringFlag{"mutex, m", "", "directory to use as a circuit instance mutex lock"},
				cli.StringFlag{"join, j", "", "join a circuit through a current member by URL"},
			},
	 	},
		// {
		// 	Name: "ls",
		// 	ShortName: "l",
		// 	//Usage: "List anchors",
		// 	Action: ls,
		// 	Flags: []cli.Flag{
		// 		cli.StringFlag{"dial, d", "", "circuit member to dial into"},
		// 	},
		// },
	}
	app.Run(os.Args)
}
