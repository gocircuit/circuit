// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package config

// AnchorConfig holds configuration parameters regarding the anchor file system server
type AnchorConfig struct {
	Addr string // Circuit address of the anchor file system worker
}
