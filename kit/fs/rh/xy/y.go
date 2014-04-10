// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package xy

import (
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

// YServer
type YServer struct {
	X circuit.PermX
}

func (y YServer) SignIn(user, dir string) (rh.Session, error) {
	r := y.X.Call("SignIn", user, dir)
	err := errors.Unpack(r[1])
	if err != nil {
		return nil, rh.Naturalize(err)
	}
	return YSession{r[0].(circuit.X)}, nil
}

func (y YServer) String() string {
	return y.X.Call("String")[0].(string)
}

// YSession
type YSession struct {
	circuit.X
}

func (y YSession) String() string {
	return y.X.Call("String")[0].(string)
}

func (y YSession) Walk(wname []string) (rh.FID, error) {
	r := y.X.Call("Walk", wname)
	err := errors.Unpack(r[1])
	if err != nil {
		return nil, rh.Naturalize(err)
	}
	return YFID{r[0].(circuit.X)}, nil
}

func (y YSession) SignOut() {
	y.X.Call("SignOut")
}

// YFID
type YFID struct {
	circuit.X
}

func (y YFID) String() string {
	return y.X.Call("String")[0].(string)
}

func (y YFID) Q() rh.Q {
	return y.X.Call("Q")[0].(rh.Q)
}

func (y YFID) Open(flag rh.Flag, intr rh.Intr) error {
	return rh.Naturalize(errors.Unpack(y.X.Call("Open", flag, newXIntr(intr))[0]))
}

func (y YFID) Create(name string, flag rh.Flag, mode rh.Mode, perm rh.Perm) (rh.FID, error) {
	r := y.X.Call("Create", name, flag, mode, perm)
	err := errors.Unpack(r[1])
	if err != nil {
		return nil, rh.Naturalize(err)
	}
	return YFID{r[0].(circuit.X)}, nil
}

func (y YFID) Clunk() error {
	r := y.X.Call("Clunk")
	err := errors.Unpack(r[0])
	if err != nil {
		return rh.Naturalize(err)
	}
	return nil
}

func (y YFID) Stat() (*rh.Dir, error) {
	r := y.X.Call("Stat")
	err := errors.Unpack(r[1])
	if err != nil {
		return nil, rh.Naturalize(err)
	}
	return r[0].(*rh.Dir), nil
}

func (y YFID) Wstat(wdir *rh.Wdir) error {
	return rh.Naturalize(errors.Unpack(y.X.Call("Wstat", wdir)[0]))
}

func (y YFID) Walk(wname []string) (rh.FID, error) {
	r := y.X.Call("Walk", wname)
	err := errors.Unpack(r[1])
	if err != nil {
		return nil, rh.Naturalize(err)
	}
	return YFID{r[0].(circuit.X)}, nil
}

func (y YFID) Read(offset int64, count int, intr rh.Intr) (rh.Chunk, error) {
	r := y.X.Call("Read", offset, count, newXIntr(intr))
	err := errors.Unpack(r[1])
	if err != nil {
		return nil, rh.Naturalize(err)
	}
	return r[0].(rh.Chunk), nil
}

func (y YFID) Write(offset int64, data rh.Chunk, intr rh.Intr) (int, error) {
	r := y.X.Call("Write", offset, data, newXIntr(intr))
	return r[0].(int), rh.Naturalize(errors.Unpack(r[1]))
}

func (y YFID) Remove() error {
	return rh.Naturalize(errors.Unpack(y.X.Call("Remove")[0]))
}

func (y YFID) Move(dir rh.FID, name string) error {
	return rh.Naturalize(errors.Unpack(y.X.Call("Move", dir, name)[0]))
}

// The presence of the FID method ensures that YFID conforms to FID at compile time.
func (y YFID) FID() rh.FID {
	return y
}
