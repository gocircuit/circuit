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

// dir represents a descendant directory within the local file system.
type dir struct {
	walk string  // absolute directory path
	f *os.File // open file for this local directory; prevents removing
}

func openDir(walk string) (*dir, error) {
	f, err := osOpenDir(walk)
	if err != nil {
		return nil, err
	}
	d := &dir{walk: walk, f: f}
	runtime.SetFinalizer(d, 
		func(d2 *dir) {
			d2.f.Close()
		},
	)
	return d, nil
}

// Path returns the local absolute path of this directory
func (d *dir) Path() string {
	return d.walk
}

// Walk returns the directory reachable from d via walk.
func (d *dir) Walk(walk ...string) (*dir, error) {
	return openDir(path.Join(append([]string{d.Path()}, walk...)...))
}

// osOpenDir opens the local absolute path name, and returns an open file handle to it, 
// as long as it is an existing directory.
func osOpenDir(name string) (*os.File, error) {
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
