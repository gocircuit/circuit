// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package rhunix

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/use/circuit"
)

func init() {
	circuit.RegisterValue(&FID{})
}

type FID struct {
	name string
	h struct {
		sync.Mutex
		root string // Local absolute path to the mount root
		walk string // File path relative to the mount root
	}
	f struct {
		sync.Mutex
		rclose   bool
		unixmode UnixMode // Cached UNIX mode for use in Q generation
		file     *os.File
	}
	d struct {
		sync.Mutex
		cached  bool      // Are directory names already cached
		entries []*rh.Dir // Cached directory entries
	}
}

func newFID(name, root, walk string, rclose bool, m os.FileMode, f *os.File) *FID {
	fid := &FID{
		name:       name,
	}
	fid.h.root = root
	fid.h.walk = walk
	fid.f.rclose = rclose
	fid.f.unixmode = UnixMode(m)
	fid.f.file = f
	//
	runtime.SetFinalizer(fid, func(x *FID) {
		x.Clunk()
	})
	// defer func() {
	// 	fid.debugf("init", nil)
	// }()
	return fid
}

func (fid *FID) rootsync() string {
	fid.h.Lock()
	defer fid.h.Unlock()
	return fid.h.root
}

func (fid *FID) walksync() string {
	fid.h.Lock()
	defer fid.h.Unlock()
	return fid.h.walk
}

func (fid *FID) unixmodesync() UnixMode {
	fid.f.Lock()
	defer fid.f.Unlock()
	return fid.f.unixmode
}

func (fid *FID) fdsync() *os.File {
	fid.f.Lock()
	defer fid.f.Unlock()
	return fid.f.file
}

func (fid *FID) Q() rh.Q {
	var id = pathQ(fid.walksync())
	fid.f.Lock()
	defer fid.f.Unlock()
	return rh.Q{
		ID:  id,
		Ver: 0,
	}
}

func (fid *FID) isdir() bool {
	fid.f.Lock()
	defer fid.f.Unlock()
	return os.FileMode(fid.f.unixmode).IsDir()
}

func (fid *FID) debugf(op string, err error) {
	f := fid.fdsync()
	var fd string = " ···· "
	if f != nil {
		fd = fmt.Sprintf("Fd=(%d)", f.Fd())
	}
	log.Printf("unixFID x=%016x %6s %s mode=%16s walk=%-20s | %v",
		unsafe.Pointer(fid), strings.ToUpper(op), fd, fid.unixmodesync().RH(), fid.walksync(), err)
}

func (fid *FID) Open(flag rh.Flag, _ rh.Intr) (err error) {
	// defer func() {
	// 	fid.debugf("open", err)
	// }()
	fid.f.Lock()
	defer fid.f.Unlock()
	if fid.f.file != nil {
		return rh.ErrBusy // file already open
	}
	//log.Printf("·····> rhflag=%v unixflag=%v", flag, RHFlag{flag}.UNIX())
	if fid.f.file, err = os.OpenFile(path.Join(fid.rootsync(), fid.walksync()), RHFlag(flag).UNIX(), 0); err != nil {
		return UnixError{err}.RH()
	}
	fid.f.rclose = flag.RemoveOnClose
	//log.Printf("=====> UNIX OPEN f=%v %s —> fd=%v", flag, path.Join(fid.rootsync(), fid.walksync()), fid.f.file.Fd())
	return nil
}

func (fid *FID) Create(name string, flag rh.Flag, mode rh.Mode, perm rh.Perm) (fid2 rh.FID, err error) {
	// defer func() {
	// 	fid.debugf("create", err)
	// }()
	name = path.Clean(name)
	if name == "." || name == ".." {
		return nil, rh.ErrClash // invalid name
	}
	if strings.Index(name, "/") >= 0 {
		return nil, rh.ErrClash // bad file name
	}
	abs := path.Join(fid.rootsync(), fid.walksync(), name)
	um := RHModePerm{mode, perm}.UNIX()
	//
	var f *os.File
	if um.IsDir() {
		if err = os.Mkdir(abs, um); err != nil {
			return nil, UnixError{err}.RH()
		}
		f, err = os.OpenFile(abs, os.O_RDONLY, 0)
	} else {
		f, err = os.OpenFile(abs, RHFlag(flag).UNIX(), um)
	}
	if err != nil {
		//log.Printf("XXXX name=%s flag=%s mode=%s perm=%s (%s)", name, flag, mode, perm, err)
		return nil, UnixError{err}.RH()
	}
	return newFID(
		path.Join(fid.name, name),
		fid.rootsync(),
		path.Join(fid.walksync(), name),
		flag.RemoveOnClose,
		um,
		f,
	), nil
}

func (fid *FID) Clunk() (err error) {
	// defer func() {
	// 	fid.debugf("clunk", nil)
	// }()
	//
	fid.f.Lock()
	var f *os.File
	f, fid.f.file = fid.f.file, nil
	var rclose bool
	rclose, fid.f.rclose = fid.f.rclose, false
	fid.f.Unlock()
	//
	if f == nil {
		return rh.ErrGone
	}
	err = f.Close()
	if rclose {
		os.Remove(path.Join(fid.rootsync(), fid.walksync()))
	}
	return UnixError{err}.RH()
}

func (fid *FID) Stat() (d *rh.Dir, err error) {
	// defer func() {
	// 	fid.debugf("stat", err)
	// }()
	//
	fid.h.Lock()
	defer fid.h.Unlock()
	//
	fid.f.Lock()
	defer fid.f.Unlock()
	//
	var um UnixMode
	// Obtain UNIX file info
	var fi os.FileInfo
	if fid.f.file == nil {
		fi, err = os.Lstat(path.Join(fid.h.root, fid.h.walk))
	} else {
		fi, err = fid.f.file.Stat()
	}
	if err != nil {
		return nil, UnixError{err}.RH()
	}
	// Obtain number of directory entries
	// XXX: This code requires an open fid. Won't work as is.
	// var numfiles int
	// if fi.IsDir() {
	// 	names, err := fid.f.file.Readdirnames(0)
	// 	if err != nil {
	// 		log.Printf("XX readdir (%s)", err)
	// 		return nil, UnixError{err}.RH()
	// 	}
	// 	numfiles = len(names)
	// }
	//
	d, um = lstat(fid.h.walk, fi, 0 /*numfiles*/)
	fid.f.unixmode = UnixMode(um) // Update cached UNIX file mode
	return d, nil
}

func (fid *FID) Wstat(wdir *rh.Wdir) (err error) {
	// defer func() {
	// 	fid.debugf("wstat", err)
	// }()
	name := path.Join(fid.rootsync(), fid.walksync())
	//
	if wdir.Perm != nil {
		var mp = RHModePerm{
			Perm: *wdir.Perm,
		}
		if err := os.Chmod(name, mp.UNIX()); err != nil {
			return UnixError{err}.RH()
		}
	}
	//
	if wdir.Mtime != nil {
		if err := os.Chtimes(name, time.Now(), *wdir.Mtime); err != nil {
			return UnixError{err}.RH()
		}
	}
	if wdir.Length != nil {
		if err := os.Truncate(name, *wdir.Length); err != nil {
			return UnixError{err}.RH()
		}
	}
	if wdir.Gid != "" {
		// Not supported yet
	}
	return nil
}

func (fid *FID) Move(dir rh.FID, name string) (err error) {
	// defer func() {
	// 	fid.debugf("move", err)
	// }()
	// Check dir is a directory
	d, err := dir.Stat()
	if err != nil {
		return err
	}
	if !d.IsDir() {
		return rh.ErrClash // not a directory
	}
	//
	fid.h.Lock() // Lock the hierarchichal structure part of the FID
	defer fid.h.Unlock()
	//
	walk := path.Join((*RHDir)(d).Walk(), name)
	if err := os.Rename(path.Join(fid.h.root, fid.h.walk), path.Join(fid.h.root, walk)); err != nil {
		return UnixError{err}.RH()
	}
	fid.h.walk = walk
	return nil
}

func (fid *FID) String() string {
	return fmt.Sprintf("unixfid·%s:(%04x)", fid.name, uint64(uintptr(unsafe.Pointer(fid)))&0xffff)
}

func (fid *FID) Walk(wname []string) (child rh.FID, err error) {
	// defer func() {
	// 	fid.debugf("walk", err)
	// }()
	child = newFID(
		path.Join(fid.name, path.Join(wname...)),
		fid.rootsync(),
		path.Join(fid.walksync(), path.Join(wname...)),
		false,
		0,
		nil,
	)
	// Calling Stat on the child initializes its m field
	if _, err := child.Stat(); err != nil {
		return nil, err
	}
	return child, nil
}

const MaxSize = 1e5 // 100KB

func (fid *FID) uf() *os.File {
	fid.f.Lock()
	defer fid.f.Unlock()
	return fid.f.file
}

func (fid *FID) Read(offset int64, count int, _ rh.Intr) (chunk rh.Chunk, err error) {
	// defer func() {
	// 	fid.debugf("read", err)
	// }()
	if fid.isdir() {
		return fid.readDir(offset, count)
	}
	return fid.readFile(offset, count)
}

func (fid *FID) readFile(offset int64, count int) (rh.ByteChunk, error) {
	buf := make([]byte, min(count, MaxSize))
	n, err := fid.uf().ReadAt(buf, offset)
	if err == io.EOF {
		err = nil // Don't report EOF error. FUSE can't understand it.
	}
	return buf[:n], UnixError{err}.RH()
}

func (fid *FID) readDir(offset int64, count int) (rh.DirChunk, error) {
	fid.d.Lock()
	defer fid.d.Unlock()
	if fid.d.cached {
		return fid.d.entries, nil
	}
	// TODO: This implementation will read all children's names on each invokation. Cache?
	root, walk := fid.rootsync(), fid.walksync()
	names, err := fid.uf().Readdirnames(0)
	if err != nil {
		return nil, UnixError{err}.RH()
	}
	if len(names) <= int(offset) {
		return nil, nil
	}
	names = names[int(offset):]
	if count > 0 && count < len(names) {
		names = names[:count]
	}
	d := make(rh.DirChunk, len(names))
	for i, n := range names {
		fi, err := os.Lstat(path.Join(root, walk, n))
		if err != nil {
			return nil, UnixError{err}.RH()
		}
		d[i], _ = lstat(walk, fi, len(names))
	}
	// Cache
	fid.d.cached, fid.d.entries = true, d
	return d, nil
}

func (fid *FID) Write(offset int64, data rh.Chunk, _ rh.Intr) (n int, err error) {
	// defer func() {
	// 	fid.debugf("write", err)
	// }()
	if fid.isdir() {
		return 0, rh.ErrClash // writing to directory
	}
	cargo := data.(rh.ByteChunk)
	cargo = cargo[:min(MaxSize, len(cargo))]
	n, err = fid.uf().WriteAt(cargo, offset)
	return n, UnixError{err}.RH()
}

func (fid *FID) Remove() (err error) {
	// defer func() {
	// 	fid.debugf("remove", err)
	// }()
	err = os.Remove(path.Join(fid.rootsync(), fid.walksync()))
	fid.f.Lock()
	fid.f.rclose = false // Disable rclose to prevent clunk from trying to remove the file
	fid.f.Unlock()
	fid.Clunk()
	return UnixError{err}.RH()
}
