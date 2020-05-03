// +build windows
// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package debug implements debugging utilities
package debug

func SavePanicTrace() {
	r := recover()
	if r == nil {
		return
	}
	panic("dumper " + r.(string))
}
