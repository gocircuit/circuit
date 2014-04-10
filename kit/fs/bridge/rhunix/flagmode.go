// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package rhunix

import (
	"fmt"
	"log"
	"os"

	"github.com/gocircuit/circuit/kit/fs/rh"
)

// RHFlag is an RH flag that has methods for conversion to UNIX flags
type RHFlag rh.Flag

func (flag RHFlag) UNIX() (unixflag int) {
	if flag.IsUnix {
		return flag.Unix
	}
	switch flag.Attr {
	case rh.Exec:
		// TODO(petar): Add permission check for execution, as described in 9P man page
		unixflag = os.O_RDONLY
	case rh.ReadWrite:
		unixflag = os.O_RDWR
	case rh.WriteOnly:
		unixflag = os.O_WRONLY
	case rh.ReadOnly:
		unixflag = os.O_RDONLY
	}
	if flag.Truncate {
		unixflag |= os.O_TRUNC
	}
	if flag.Create {
		unixflag |= os.O_CREATE
	}
	return
}

// UnixFlag is a UNIX flag to Open with methods for conversion to an RH flag
type UnixFlag int

// RH returns the corresponding UNIX flag.
// Flags O_APPEND, O_CREATE, O_EXCL, O_SYNC are ignored.
func (uf UnixFlag) RH() (rhf rh.Flag) {
	u := (int)(uf)
	switch u & 3 {
	case os.O_RDWR:
		rhf.Attr = rh.ReadWrite
	case os.O_WRONLY:
		rhf.Attr = rh.WriteOnly
	//case os.O_TRUNC, os.O_CREATE:
	//	rhf.Attr = rh.WriteOnly
	case os.O_APPEND: // Append write, random read
		rhf.Attr = rh.ReadWrite
	case os.O_RDONLY: // O_RDONLY == 0
		rhf.Attr = rh.ReadOnly
	}
	if os.O_TRUNC&u != 0 {
		rhf.Truncate = true
	}
	if os.O_CREATE&u != 0 {
		rhf.Create = true
	}
	rhf.IsUnix = true
	rhf.Unix = u
	return
}

// RHModePerm is an RH mode/perm pair with methods for conversion to a UNIX mode
type RHModePerm struct {
	rh.Mode
	rh.Perm
}

func (mp RHModePerm) String() string {
	return fmt.Sprintf("[%5s %s]", mp.Mode, mp.Perm)
}

func (mp RHModePerm) UNIX() (umode os.FileMode) {
	if mp.Mode.IsUnix {
		return mp.Mode.Unix
	}
	switch mp.Mode.Attr {
	case rh.ModeFile, rh.ModeIO:
		// IO channels are shown as regular files to FUSE.
	case rh.ModeDir:
		umode = os.ModeDir
	case rh.ModeLog:
		umode = os.ModeAppend
	case rh.ModeMutex:
		umode = os.ModeExclusive
	case rh.ModeRef:
		log.Println("symlinks not supported")
		umode = os.ModeSymlink
	case rh.ModeUnknown:
		umode = os.ModeSocket // Not supported. This UNIX mode will cause follow up ops to fail, as intended.
	}
	umode |= os.FileMode(mp.Perm) // Copy permission
	return
}

// UnixMode is a UNIX file mode with method for conversion to an RH mode/perm pair and QType
type UnixMode os.FileMode

func (mode UnixMode) RH() (mp RHModePerm) {
	m := os.FileMode(mode)
	mp.Mode.Attr = rh.ModeFile
	if m&os.ModeDir != 0 {
		mp.Mode.Attr = rh.ModeDir
	}
	if m&os.ModeAppend != 0 {
		mp.Mode.Attr = rh.ModeLog
	}
	if m&os.ModeExclusive != 0 {
		mp.Mode.Attr = rh.ModeMutex
	}
	if m&os.ModeSymlink != 0 {
		mp.Mode.Attr = rh.ModeRef
	}
	if m&os.ModeDevice != 0 {
		mp.Mode.Attr = rh.ModeUnknown
	}
	if m&os.ModeNamedPipe != 0 {
		mp.Mode.Attr = rh.ModeUnknown
	}
	if m&os.ModeSocket != 0 {
		mp.Mode.Attr = rh.ModeUnknown
	}
	mp.Mode.IsUnix = true
	mp.Mode.Unix = os.FileMode(mode)
	// ModeCharDevice and ModeSticky have no counterparts
	mp.Perm = rh.Perm(m & os.ModePerm)
	return
}
