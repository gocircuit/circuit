// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package kill has the side effect of installing a KILL signal handler that throws a panic
package kill

import "github.com/hoijui/circuit/pkg/kit/debug"

func init() {
	debug.InstallKillPanic()
}
