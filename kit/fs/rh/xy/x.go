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

func init() {
	circuit.RegisterValue(XServer{})
	circuit.RegisterValue(XSession{})
	circuit.RegisterValue(XFID{})
	circuit.RegisterValue(&XIntr{})
}

// XServer
type XServer struct {
	rh.Server
}

func (x XServer) SignIn(user, dir string) (circuit.X, error) {
	fsys, err := x.Server.SignIn(user, dir)
	if err != nil {
		return nil, errors.Pack(err)
	}
	return circuit.PermRef(XSession{fsys}), nil
}

// XSession
type XSession struct {
	rh.Session
}

func (x XSession) Walk(name []string) (circuit.X, error) {
	fid, err := x.Session.Walk(name)
	if err != nil {
		return nil, errors.Pack(err)
	}
	return circuit.Ref(XFID{fid}), nil
}

// XFID
type XFID struct {
	rh.FID
}

func (x XFID) Open(flag rh.Flag, xintr circuit.X) error {
	yintr := NewNoPanicYIntr(xintr)
	defer yintr.Clunk()
	return x.FID.Open(flag, yintr.Intr)
}

func (x XFID) Clunk() error {
	return errors.Pack(x.FID.Clunk())
}

func (x XFID) Create(name string, flag rh.Flag, mode rh.Mode, perm rh.Perm) (circuit.X, error) {
	fid, err := x.FID.Create(name, flag, mode, perm)
	if err != nil {
		return nil, errors.Pack(err)
	}
	return circuit.Ref(XFID{fid}), nil
}

func (x XFID) Walk(wname []string) (circuit.X, error) {
	fid, err := x.FID.Walk(wname)
	if err != nil {
		return nil, errors.Pack(err)
	}
	return circuit.Ref(XFID{fid}), nil
}

func (x XFID) Move(dir rh.FID, name string) error {
	return x.FID.Move(dir, name)
}

func (x XFID) Read(offset int64, count int, xintr circuit.X) (rh.Chunk, error) {
	yintr := NewNoPanicYIntr(xintr)
	defer yintr.Clunk()
	return x.FID.Read(offset, count, yintr.Intr)
}

func (x XFID) Write(offset int64, data rh.Chunk, xintr circuit.X) (int, error) {
	yintr := NewNoPanicYIntr(xintr)
	defer yintr.Clunk()
	return x.FID.Write(offset, data, yintr.Intr)
}

// XIntr
type XIntr struct {
	intr rh.Intr
}

func newXIntr(intr rh.Intr) circuit.X {
	return circuit.Ref(&XIntr{
		intr: intr,
	})
}

func (x *XIntr) Wait() rh.Prompt {
	return <-x.intr
}
