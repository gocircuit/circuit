// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package anchorfs

import (
	"github.com/hoijui/circuit/pkg/kit/module"
)

var mod = module.Slot{Name: "anchor file system"}

// Bind is used internally to bind an implementation of this package to the public methods of this package
func Bind(v System) {
	mod.Set(v)
}

func get() System {
	return mod.Get().(System)
}

// Created returns a slive of anchor directories within which this worker has created files with CreateFile.
func Created() []string {
	return get().Created()
}

// OpenDir opens the anchor directory anchor
func OpenDir(anchor string) (Dir, error) {
	return get().OpenDir(anchor)
}

// OpenFile opens the anchor file anchor
func OpenFile(anchor string) (File, error) {
	return get().OpenFile(anchor)
}
