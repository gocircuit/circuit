// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package load supplies an implementation of the anchorfs circuit module
package load

import (
	"github.com/gocircuit/circuit/sys/anchor"
	"github.com/gocircuit/circuit/sys/anchor/xy"
	use "github.com/gocircuit/circuit/use/anchorfs"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
)

var anchorKeepalive xy.YKeepalive

func Load(srvAddr n.Addr, workerAddr n.Addr, anchors []string) {
	if anchorKeepalive != nil {
		panic("connecting to anchor file system twice")
	}
	x := circuit.Dial(srvAddr, "anchor")
	y := xy.YSys{x}
	var err error
	if anchorKeepalive, err = y.Create(anchors, workerAddr); err != nil {
		panic(err)
	}
	use.Bind(&sys{y: y, created: anchors})
}

// sys
type sys struct {
	y       xy.YSys
	created []string
}

func (s *sys) OpenFile(fullpath string) (use.File, error) {
	f, err := s.y.OpenFile(fullpath)
	if err != nil {
		return nil, err
	}
	return file{f}, nil
}

func (s *sys) OpenDir(fullpath string) (use.Dir, error) {
	d, err := s.y.OpenDir(fullpath)
	if err != nil {
		return nil, err
	}
	return dir{d}, nil
}

func (s *sys) Created() []string {
	return s.created
}

// dir
type dir struct {
	xy.YDir
}

func (d dir) OpenFile(fullpath string) (use.File, error) {
	f, err := d.YDir.OpenFile(fullpath)
	if err != nil {
		return nil, err
	}
	return file{f}, nil
}

func (d dir) OpenDir(fullpath string) (use.Dir, error) {
	b, err := d.YDir.OpenDir(fullpath)
	if err != nil {
		return nil, err
	}
	return dir{b}, nil
}

// file
type file struct {
	*anchor.File
}

func (f file) Path() string {
	return f.File.Path()
}

func (f file) Anchor() n.Addr {
	return f.File.Anchor().(n.Addr)
}
