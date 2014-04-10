// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"runtime"
	"time"

	use "github.com/gocircuit/circuit/use/anchorfs"
)

// Dir is a directory handle.
// It points to a directory Node and prevents it from being reclaimed by the file system's garbage collection mechanism.
type Dir struct {
	h *Handle
}

func newDir(h *Handle) *Dir {
	d := &Dir{h: h}
	runtime.SetFinalizer(d, func(x *Dir) {
		x.h.Close()
	})
	return d
}

func (d *Dir) node() *Node {
	return d.h.Node
}

// Path returns the fully-qualified path of this directory.
func (d *Dir) Path() string {
	return d.h.Path
}

func (d *Dir) List() (rev use.Rev, files, dirs []string) {
	fnodes, dnodes, rev := d.h.List()
	files, dirs = flatten(fnodes, dnodes)
	return rev, files, dirs
}

func flatten(fnodes, dnodes []*Handle) (files, dirs []string) {
	files, dirs = make([]string, len(fnodes)), make([]string, len(dnodes))
	for i, fn := range fnodes {
		files[i] = fn.Node.Path
	}
	for i, dn := range dnodes {
		dirs[i] = dn.Node.Path
	}
	return
}

func (d *Dir) Change(sinceRev use.Rev) (rev use.Rev, files, dirs []string) {
	fnodes, dnodes, rev := d.h.Change(sinceRev)
	files, dirs = flatten(fnodes, dnodes)
	return rev, files, dirs
}

func (d *Dir) ChangeExpire(sinceRev use.Rev, expire time.Duration) (rev use.Rev, files, dirs []string, err error) {
	fnodes, dnodes, rev, err := d.h.ChangeExpire(sinceRev, expire)
	files, dirs = flatten(fnodes, dnodes)
	return rev, files, dirs, err
}

func (d *Dir) OpenDir(name string) (*Dir, error) {
	h := d.h.Open(name)
	if h == nil {
		return nil, ErrExist
	}
	if h.IsFile() {
		return nil, ErrKind
	}
	return newDir(h), nil
}

func (d *Dir) OpenFile(name string) (*File, error) {
	n := d.h.Open(name)
	if n == nil {
		return nil, ErrExist
	}
	if !n.IsFile() {
		return nil, ErrKind
	}
	return newFile(n), nil
}
