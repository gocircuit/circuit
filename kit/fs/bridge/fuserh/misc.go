// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package fuserh

import (
	"fmt"
	"io"
	"log"
	"path"
	"strings"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/fuse"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

func Debugf(format string, arg ...interface{}) {
	log.Printf("FUSE/RH "+format, arg...)
}

func Debugln(arg ...interface{}) {
	log.Printf("FUSE/RH %s", fmt.Sprintln(arg...))
}

// IsJunk
func IsJunk(name string) bool {
	for _, j := range JunkFiles {
		if strings.HasPrefix(name, j) {
			return true
		}
	}
	return false
}

func IsJunkRequest(req fuse.Request) bool {
	switch q := req.(type) {
	case *fuse.LookupRequest:
		return IsJunk(q.Name)
	case *fuse.StatfsRequest:
		return true
	}
	return false
}

// RHError
type RHError struct {
	error
}

// FUSE converts errors returned from an RH interface to FUSE errors
func (rhe RHError) FUSE() fuse.Error {
	if rhe.error == nil {
		return nil
	}
	//log.Printf("Error RH->FUSE: %#v", rhe.error)
	re, ok := rhe.error.(rh.Error)
	if !ok {
		panic("u")
	}
	switch {
	case re.IsEqual(rh.ErrExist):
		return fuse.EACCES
	case re.IsEqual(rh.ErrNotExist):
		return fuse.ENOENT
	case re.IsEqual(rh.ErrBusy):
		return fuse.EBUSY
	case re.IsEqual(rh.ErrGone):
		return fuse.EBADF
	case re.IsEqual(rh.ErrPerm):
		return fuse.EPERM
	case re.IsEqual(rh.ErrIO):
		return fuse.EIO
	case re.IsEqual(rh.ErrEOF):
		return nil // No FUSE equivalent for EOF, right? Return no error for now.
	case re.IsEqual(rh.ErrClash):
		return fuse.EINVAL
	case re.IsEqual(rh.ErrIntr):
		return fuse.EINTR
	}
	return fuse.ESTALE
}

//
func split(s string) []string {
	w := strings.Split(path.Clean(s), "/")
	if len(w) > 0 && w[0] == "" {
		w = w[1:]
	}
	if len(w) > 0 && w[len(w)-1] == "" {
		w = w[:len(w)-1]
	}
	return w
}

//
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// EOF
type EOF struct {
	r <-chan error
	sync.Mutex
	w chan<- error
}

func (eof *EOF) Init() {
	ch := make(chan error, 1)
	eof.r, eof.w = ch, ch
}

func (eof *EOF) EOF() error {
	err, ok := <-eof.r
	if !ok {
		return io.ErrClosedPipe
	}
	return err
}

func (eof *EOF) Close(err error) {
	eof.Lock()
	defer eof.Unlock()
	if eof.w == nil {
		return
	}
	eof.w <- err
	close(eof.w)
	eof.w = nil
}

func (eof *EOF) Chan() <-chan error {
	return eof.r
}
