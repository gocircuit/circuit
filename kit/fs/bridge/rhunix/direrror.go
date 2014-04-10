// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package rhunix

import (
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"syscall"

	"github.com/gocircuit/circuit/kit/fs/rh"
)

func pathQ(p string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(p))
	return h.Sum64()
}

// RHDir is an RH directories with methods for obtaining various corresponding UNIX statistics
type RHDir rh.Dir

func (d *RHDir) Walk() string {
	return (*rh.Dir)(d).Aux.(*DirAux).Walk
}

func (d *RHDir) UnixLen() (int64, bool) {
	rhd := (*rh.Dir)(d)
	if rhd.Aux == nil {
		return 0, false
	}
	return rhd.Aux.(*DirAux).UnixLen, true
}

func (d *RHDir) Inode() uint64 {
	const RootInodeID = 1
	if d.Name == "/" {
		return RootInodeID
	}
	return d.Q.Hash64()
}

// DirAux is auxiliary data that this file system attaches to rh.Dir structures
type DirAux struct {
	Walk    string
	UnixLen int64
}

func init() {
	gob.Register(&DirAux{})
}

func (x *DirAux) String() string {
	return x.Walk
}

// unwrap
func unwrap(err error) error {
	switch pe := err.(type) {
	case nil:
		return nil
	case *syscall.Errno:
		return pe
	case *os.PathError:
		return pe.Err
	case *os.LinkError:
		return pe.Err
	case *os.SyscallError:
		return pe.Err
	}
	switch err {
	case io.EOF, io.ErrClosedPipe, io.ErrNoProgress, io.ErrShortBuffer, io.ErrShortWrite, io.ErrUnexpectedEOF:
		return err
	case os.ErrInvalid:
		return err
	}
	panic(fmt.Sprintf("unknown system error %#v %#T", err, err))
}

// UnixError
type UnixError struct {
	error
}

func (ue UnixError) RH() error {
	if ue.error == nil {
		return nil
	}
	if _, ok := ue.error.(rh.Error); ok {
		panic(fmt.Sprintf("expecting rh.Error, got %#v %#T", ue.error, ue.error))
	}
	if os.IsExist(ue.error) {
		return rh.ErrExist
	}
	if os.IsNotExist(ue.error) {
		return rh.ErrNotExist
	}
	if os.IsPermission(ue.error) {
		return rh.ErrPerm
	}
	//
	switch unwrap(ue.error) {
	case os.ErrInvalid, syscall.EINVAL:
		return rh.ErrClash
	case syscall.EBUSY:
		return rh.ErrBusy
	case syscall.EBADF:
		return rh.ErrGone
	case syscall.ESTALE:
		return rh.ErrGone
	case syscall.EIO:
		return rh.ErrIO
	case syscall.EACCES:
		return rh.ErrExist
	case io.EOF:
		return rh.ErrEOF
	case io.ErrClosedPipe, io.ErrUnexpectedEOF:
		return rh.ErrIO
	case io.ErrNoProgress, io.ErrShortBuffer, io.ErrShortWrite:
		return rh.ErrIO
	case syscall.EISDIR:
		return rh.ErrClash
	case os.ErrInvalid:
		return rh.ErrClash
	}
	//
	log.Printf("unknown sys error (%s)", unwrap(ue.error))
	return rh.ErrClash
}

// lstat
func lstat(walk string, fi os.FileInfo, numfiles int) (*rh.Dir, UnixMode) {
	uid, gid, atime, ver := readSysStat(fi.Sys())
	//
	um := UnixMode(fi.Mode())
	mp := um.RH()
	return &rh.Dir{
		Q: rh.Q{
			ID:  pathQ(walk),
			Ver: ver,
		},
		Mode:   mp.Mode,
		Perm:   mp.Perm,
		Atime:  atime,
		Mtime:  fi.ModTime(),
		Length: int64(numfiles), // Number of child files, by RH spec
		Name:   fi.Name(),
		Uid:    uid,
		Gid:    gid,
		Aux: &DirAux{
			Walk:    walk,
			UnixLen: fi.Size(),
		},
	}, um
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
