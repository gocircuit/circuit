// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package fuserh

import (
	"io"

	"github.com/gocircuit/circuit/kit/fs/fuse"
)

func (r *RH) init() interface{} {
	return &fuse.InitResponse{
		MaxWrite: 60e3, // 60K IP packet size  // uint32(syscall.Getpagesize()),
	}
}

func (r *RH) statfs() interface{} {
	return &fuse.StatfsResponse{
		Blocks:  0,    // Total data blocks in file system.
		Bfree:   0,    // Free blocks in file system.
		Bavail:  0,    // Free blocks in file system if you're not root.
		Files:   0,    // Total files in file system.
		Ffree:   0,    // Free files in file system.
		Bsize:   8192, // Block size
		Namelen: 4096, // Maximum file name length?
		Frsize:  0,    // ?
	}
}

func (r *RH) destroy() interface{} {
	r.ssn.SignOut()
	r.fuseEOF.Close(io.EOF)
	return nil
}
