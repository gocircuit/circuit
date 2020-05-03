// +build !windows
// +build arm arm64

// Copyright 2019 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.

// Package debug implements debugging utilities
package debug

import (
	"syscall"
)

func Dup2(oldfd int, newfd int) {
	syscall.Dup3(oldfd, newfd, 0)
}
