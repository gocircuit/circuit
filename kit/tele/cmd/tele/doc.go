// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

//  Package teleport implements a command-line tool for utilizing teleport transport between legacy clients and servers.
package main

import (
	"fmt"
	"os"
)

// TODO: Add limit on concurrent connections.

const help = `
Teleport Transport Tool
Part of The Go Circuit Project 2013, http://gocircuit.org
______________________________________________________________________________________
Operating diagram:

                                     client/input
                                           |
        +---------------+                  | +-------------+
        | USER's CLIENT +---- localhost -->• | TELE CLIENT +-----+
        +---------------+                    +-------------+     |
                                                                 ≈
                                                                 |
       CLIENT-SIDE                                               |
  ·····················································  UNRELIABLE NETWORK  ·····
       SERVER-SIDE                                               |
                                                                 |
                                                                 ≈
        +---------------+                    +-------------+     |
        | USER's SERVER | •<-- localhost ----+ TELE SERVER | •<--+
        +---------------+ |                  +-------------+ |
                          |                                  |
      	          server/input                server+client/output

______________________________________________________________________________________
TELEPORT SERVER:

tele -server -in=input_addr [-out=output_addr]

In server regime, the teleport tool will accept TELEPORT connections incoming
to output_addr and forward/proxy them to the TCP server listening on input_addr.
If output_addr is not specified, the tool will use an available port and print it.

______________________________________________________________________________________
TELEPORT CLIENT:

tele -client [-in=input_address] -out=output_addr

In client regime, the teleport tool will accept TCP connections incoming
to input_addr and forward/proxy them to the TELEPORT server listening on output_addr.
If input_addr is not specified, the tool will use an available port and print it.

`

/*
______________________________________________________________________________________
COMMON OPTIONS:

-max=M  Limit the number of concurrently open connections to M.
*/

func usage() {
	fatalf(help)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
