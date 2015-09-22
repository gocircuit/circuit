// +build !windows
// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package os provides application-level utilities that are implemented using OS facilities (like lock files)
package lockfile

import (
	"os"
	"syscall"
)

// LockFile represents an exclusive, advisory OS file lock
type LockFile struct {
	file *os.File
}

// Create creates a file with an exclusive, advisory lock
func Create(name string) (*LockFile, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		return nil, err
	}
	return &LockFile{file}, nil
}

// Release releases the OS file lock and closes the respective file
func (lf *LockFile) Release() error {
	if err := syscall.Flock(int(lf.file.Fd()), syscall.LOCK_UN); err != nil {
		lf.file.Close()
		lf.file = nil
		return err
	}
	if err := lf.file.Close(); err != nil {
		lf.file = nil
		return err
	}
	return nil
}
