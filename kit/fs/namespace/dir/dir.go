// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package dir

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/gocircuit/circuit/kit/fs/rh"
)

type Dir struct {
	rh.ReadOnlyFID
	mtime time.Time
	sync.Mutex
	atime    time.Time
	children map[string]*childFID
	writable bool
}

func NewDir() *Dir {
	now := time.Now()
	return &Dir{
		mtime: now,
		atime: now,
		children:  make(map[string]*childFID),
	}
}

// User interface

func (d *Dir) AllowCreate() {
	d.Lock()
	defer d.Unlock()
	d.writable = true
}

func (d *Dir) DisallowCreate() {
	d.Lock()
	defer d.Unlock()
	d.writable = true
}

func (d *Dir) Len() int {
	d.Lock()
	defer d.Unlock()
	return len(d.children)
}

func (d *Dir) AddChild(name string, child rh.FID) (rh.FID, error) {
	d.Lock()
	defer d.Unlock()
	if _, ok := d.children[name]; ok {
		return nil, rh.ErrExist
	}
	d.children[name] = NewChildFID(d, name, child)
	return d.children[name], nil
}

func (d *Dir) Clear() {
	d.Lock()
	defer d.Unlock()
	d.children = make(map[string]*childFID)
}

func (d *Dir) RemoveChild(name string) {
	d.Lock()
	defer d.Unlock()
	delete(d.children, name)
}

func (d *Dir) FID() rh.FID {
	return (*DirFID)(d)
}

// RH-interface

type DirFID Dir

func (fid *DirFID) String() string {
	return fmt.Sprintf("dir(%02x)", uint64(uintptr(unsafe.Pointer(fid)))&0xff)
}

func (fid *DirFID) Q() rh.Q {
	return rh.Q{
		ID:  interfaceHash(fid), // QID is a hash of the pointer value
		Ver: 1,
	}
}

func (fid *DirFID) Open(flag rh.Flag, _ rh.Intr) error {
	return nil
}

func (fid *DirFID) Clunk() error {
	return nil
}

func (fid *DirFID) Stat() (*rh.Dir, error) {
	fid.Lock()
	defer fid.Unlock()
	//
	d := &rh.Dir{
		Q: fid.Q(),
		Mode: rh.Mode{
			Attr:     rh.ModeDir,
			IsHidden: false,
		},
		Atime:  fid.atime,
		Mtime:  fid.mtime,
		Name:   "¡?", // Name is set by the wrapping owner
		Length: int64(len(fid.children)),
		Uid:    rh.UID(),
		Gid:    rh.GID(),
	}
	if fid.writable {
		d.Perm = 0777 // rwxrwxrwx
	} else {
		d.Perm = 0555 // r-xr-xr-x
	}
	return d, nil
}

func (fid *DirFID) Walk(wname []string) (rh.FID, error) {
	fid.Lock()
	defer fid.Unlock()
	//
	if len(wname) == 0 {
		return fid, nil
	}
	f, ok := fid.children[wname[0]]
	if !ok {
		return nil, rh.ErrNotExist
	}
	return f.Walk(wname[1:])
}

func (fid *DirFID) Read(offset int64, count int, _ rh.Intr) (rh.Chunk, error) {
	fid.Lock()
	defer fid.Unlock()
	//
	r := make(rh.DirChunk, 0, len(fid.children))
	for _, dir := range fid.children {
		d, _ := dir.Stat()
		r = append(r, d)
	}
	return r, nil
}

func (fid *DirFID) Remove() error {
	fid.Lock()
	defer fid.Unlock()
	if len(fid.children) > 0 {
		return rh.ErrPerm
	}
	return nil
}
