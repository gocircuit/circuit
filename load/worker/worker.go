// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Importing package worker has the side effect of turning your program into a circuit worker executable
package worker

import (
	_ "github.com/gocircuit/circuit/kit/debug/kill"
	_ "github.com/gocircuit/circuit/load/cmd"
)

func init() {
	// After package load installs and activates all circuit-related logic,
	// this function blocks forever, never allowing execution of main.
	<-(chan struct{})(nil)
}
