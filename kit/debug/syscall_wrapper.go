// +build !windows,!arm
// Copyright 2019 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.

// Package debug implements debugging utilities
package debug

import (
	"syscall"
)

func Dup2(oldfd int, newfd int) {
	syscall.Dup2(oldfd, newfd)
}
