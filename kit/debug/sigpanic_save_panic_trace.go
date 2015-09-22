// +build !windows
// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package debug implements debugging utilities
package debug

import (
	"fmt"
	"os"
	"syscall"
)

func SavePanicTrace() {
	r := recover()
	if r == nil {
		return
	}
	// Redirect stderr
	file, err := os.Create("panic")
	if err != nil {
		panic("dumper (no file) " + r.(fmt.Stringer).String())
	}
	syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	// TRY: defer func() { file.Close() }()
	panic("dumper " + r.(string))
}
