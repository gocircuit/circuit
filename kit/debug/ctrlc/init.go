// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package ctrlc has the side effect of installing a Ctrl-C signal handler that throws a panic
package ctrlc

import "github.com/gocircuit/circuit/kit/debug"

func init() {
	debug.InstallCtrlCPanic()
}
