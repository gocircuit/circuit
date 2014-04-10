// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package fuse

import (
	"syscall"
)

// An Error is a FUSE error.
type Error interface {
	errno() int32
}

const (
	// ENOSYS indicates that the call is not supported.
	ENOSYS = Errno(syscall.ENOSYS) // n/a

	// ESTALE is used by Serve to respond to violations of the FUSE protocol.
	ESTALE = Errno(syscall.ESTALE) // clash

	EACCES = Errno(syscall.EACCES) // exist
	ENOENT = Errno(syscall.ENOENT) // not exist

	EBUSY = Errno(syscall.EBUSY) // busy
	EBADF = Errno(syscall.EBADF) // gone

	EIO   = Errno(syscall.EIO)   // io
	EPERM = Errno(syscall.EPERM) // perm

	EINTR  = Errno(syscall.EINTR)  // intr
	EINVAL = Errno(syscall.EINVAL) // clash

	// ETIMEDOUT
	// E2BIG
	// ENODEV
	// EAGAIN
	// ENOMEM
	// EISDIR, ENOTDIR
	// EPIPE
	// EROFS
)

// Errno implements Error using a syscall.Errno.
type Errno syscall.Errno

func (e Errno) errno() int32 {
	return int32(e)
}

func (e Errno) String() string {
	return syscall.Errno(e).Error()
}
