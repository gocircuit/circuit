// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package fuserh

import (
	"os"
	"strconv"
	"time"

	"github.com/gocircuit/circuit/kit/fs/fuse"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/bridge/rhunix"
)

// accessMask, see man page for access
type accessMask uint32

func (am accessMask) RH() rh.Flag {
	panic("u/i")
}

// fuseOpenFlag
type fuseOpenFlag uint32

// Return the Go's os package equivalent of the FUSE open request flags
func (f fuseOpenFlag) Unix() int {
	return int(f)
}

// fuseAttr
func fuseAttr(dir *rh.Dir) (attr fuse.Attr, attrValid, entryValid time.Duration) {
	uid, err := strconv.Atoi(dir.Uid)
	if err != nil {
		uid = os.Getuid()
	}
	gid, err := strconv.Atoi(dir.Gid)
	if err != nil {
		gid = os.Getgid()
	}
	unixlen, ok := (*rhunix.RHDir)(dir).UnixLen()
	if !ok {
		unixlen = dir.Length
	}
	return fuse.Attr{
		Inode:  (*rhunix.RHDir)(dir).Inode(),
		Size:   uint64(unixlen),
		Blocks: (uint64(dir.Length) + 8191) / 8192, // size in blocks
		Atime:  dir.Atime,
		Mtime:  dir.Mtime,
		Ctime:  dir.Mtime, // time of last inode change; hacked
		Crtime: dir.Mtime, // time of creation (OS X only)
		Mode: rhunix.RHModePerm{
			Mode: dir.Mode,
			Perm: dir.Perm,
		}.UNIX(), // file mode
		Nlink: 1, // number of links; rsc says, works for directories! - see FUSE FAQ
		Uid:   uint32(uid),
		Gid:   uint32(gid),
		Rdev:  0, // device numbers
		Flags: 0, // chflags(2) flags (OS X only)
	}, 0, 0
}

// sattrWdir converts a fuse.SetattrRequest into a rh.Wdir
func sattrWdir(sattr *fuse.SetattrRequest) (wdir *rh.Wdir) {
	//
	var valid = sattr.Valid
	switch {
	case valid.Uid(), valid.Atime():
		// ignore
	}
	//
	wdir = &rh.Wdir{}
	if valid.Mode() {
		var mp = rhunix.UnixMode(sattr.Mode).RH()
		wdir.Perm = &mp.Perm
	}
	if valid.Mtime() {
		var mtime = sattr.Mtime
		wdir.Mtime = &mtime
	}
	if valid.Gid() {
		wdir.Gid = strconv.Itoa(int(sattr.Gid))
	}
	if valid.Size() {
		var l = int64(sattr.Size)
		wdir.Length = &l
	}
	return
}
