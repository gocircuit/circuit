// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package config

import "github.com/gocircuit/circuit/use/n"

// SparkConfig captures a few worker startup parameters that can be configured on each execution
type SparkConfig struct {
	// ID is the ID of the worker instance
	ID n.WorkerID

	// BindAddr is the network address the worker will listen to for incoming connections
	BindAddr string

	// Host is the host name of the hosting machine
	Host string

	// Anchor is the set of anchor directories that the worker registers with
	Anchor []string
}

// DefaultSpark is the default configuration used for workers started from the command line, which
// are often not intended to be contacted back from other workers
var DefaultSpark *SparkConfig
