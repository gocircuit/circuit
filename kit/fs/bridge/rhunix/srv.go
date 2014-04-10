// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package rhunix is an RH interface to the local UNIX file system.
//
// Formally, it is a ––>RH––>UNIX.
package rhunix

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"

	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/use/circuit"
)

// Prior art:
//
//	http://plan9.bell-labs.com/magic/man2html/4/u9fs
//	https://www.usenix.org/legacy/events/usenix05/tech/freenix/full_papers/hensbergen/hensbergen.pdf
//

func init() {
	circuit.RegisterValue(&Server{})
	circuit.RegisterValue(&Session{})
}

// Server is an RH server that can attach subtrees of a local UNIX directory.
type Server struct {
	name string
	root string
}

// New creates a new RH server for the local directory localroot, named name.
func New(name, localroot string) (*Server, error) {
	abs, err := filepath.Abs(localroot)
	if err != nil {
		return nil, err
	}
	return &Server{
		name: name,
		root: abs,
	}, nil
}

// SignIn creates a new file system server for the user, rooted at root within the server's local root directory.
func (srv *Server) SignIn(user, root string) (ssn rh.Session, err error) {

	// Slash
	slash := newFID("/", path.Join(srv.root, root), "/", false, 0, nil)
	dir, err := slash.Stat() // Calling Stat will initiailze fid, by updating its mode cache
	if err != nil {
		return nil, err
	}
	if !dir.IsDir() {
		return nil, rh.ErrNotExist // attaching non-directory
	}
	//
	ssn_ := &Session{
		user:  user,
		root:  root,
		slash: slash,
	}
	ssn = ssn_
	ssn_.name = fmt.Sprintf("%s:(%02x)", srv.String(), uint64(uintptr(unsafe.Pointer(ssn_)))&0xff)
	runtime.SetFinalizer(ssn, func(x *Session) {
		x.SignOut()
	})
	return ssn, nil
}

func (srv *Server) String() string {
	return srv.name
}

type Session struct {
	name string
	user string // User attached to this ssn
	root string // Subtree path to this file system relative to the server root
	sync.Mutex
	slash *FID // FID of the file system root directory
}

func (ssn *Session) Walk(wname []string) (fid rh.FID, err error) {
	return ssn.slashFID().Walk(wname)
}

func (ssn *Session) slashFID() *FID {
	ssn.Lock()
	defer ssn.Unlock()
	return ssn.slash
}

func (ssn *Session) String() string {
	return ssn.name
}

func (ssn *Session) SignOut() {
	// For the moment, we don't clunk open FIDs on detach. Should do eventually.
	ssn.Lock()
	defer ssn.Unlock()
	ssn.slash = nil
}
