// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"errors"
	"os"
	"path"
	"runtime"
)

// Dir represents a descendant directory within the local circuit-mounted file system.
type Dir struct {
	walk string  // absolute directory path
	dir *os.File // open file for this local directory; prevents removing
}

func OpenDir(walk string) (*Dir, error) {
	f, err := openDir(walk)
	if err != nil {
		return nil, err
	}
	d := &Dir{walk: walk, dir: f}
	runtime.SetFinalizer(d, 
		func(d2 *Dir) {
			d.dir.Close()
		},
	)
	return d, nil
}

// Path returns the local absolute path of this directory
func (d *Dir) Path() string {
	return d.walk
}

// Walk returns the directory reachable from d via walk.
func (d *Dir) Walk(walk ...string) (*Dir, error) {
	return OpenDir(path.Join(append([]string{d.Path()}, walk...)...))
}

// openDir opens the local absolute path name, and returns an open file handle to it, 
// as long as it is an existing directory.
func openDir(name string) (*os.File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if !fi.IsDir() {
		f.Close()
		return nil, errors.New("not a directory")
	}
	return f, nil
}
