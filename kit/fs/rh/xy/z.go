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
	"sync"
)

// NoPanicYServer
type NoPanicYServer struct {
	X circuit.PermX
}

func (y NoPanicYServer) SignIn(user, dir string) (ssn rh.Session, err error) {
	defer func() {
		if r := recover(); r != nil {
			ssn, err = nil, rh.ErrIO // errors.NewPanic(r)
		}
	}()
	//
	r := y.X.Call("SignIn", user, dir)
	if err = errors.Unpack(r[1]); err != nil {
		return nil, rh.Naturalize(err)
	}
	return &NoPanicYSession{X: r[0].(circuit.X)}, nil
}

func (y NoPanicYServer) String() string {
	return y.X.Call("String")[0].(string)
}

// NoPanicYSession
type NoPanicYSession struct {
	circuit.X
	c__     sync.Mutex
	c__name string
}

func (y *NoPanicYSession) String() (s string) {
	y.c__.Lock()
	defer y.c__.Unlock()
	if y.c__name != "" {
		return y.c__name
	}
	defer func() {
		if r := recover(); r != nil {
			s = "¡xerr!"
			y.c__name = s
		}
	}()
	y.c__name = y.X.Call("String")[0].(string)
	return y.c__name
}

func (y *NoPanicYSession) Walk(wname []string) (fid rh.FID, err error) {
	//log.Println(y.String(), "Walk")
	defer func() {
		if r := recover(); r != nil {
			fid, err = nil, rh.ErrIO //errors.NewPanic(r)
		}
	}()
	//
	r := y.X.Call("Walk", wname)
	if err = errors.Unpack(r[1]); err != nil {
		return nil, rh.Naturalize(err)
	}
	return &NoPanicYFID{X: r[0].(circuit.X)}, nil
}

func (y *NoPanicYSession) SignOut() {
	//log.Println(y.String(), "SignOut")
	defer func() {
		recover()
	}()
	//
	y.X.Call("SignOut")
}

// NoPanicYFID
type NoPanicYFID struct {
	circuit.X // Non-perm pointer, so FIDs can be garbage-collected
	c__       sync.Mutex
	c__name   string
}

func (y *NoPanicYFID) String() (s string) {
	y.c__.Lock()
	defer y.c__.Unlock()
	if y.c__name != "" {
		return y.c__name
	}
	defer func() {
		if r := recover(); r != nil {
			s = "¡xerr!"
			y.c__name = s
		}
	}()
	y.c__name = y.X.Call("String")[0].(string)
	return y.c__name
}

func (y *NoPanicYFID) Q() rh.Q {
	defer func() {
		recover()
	}()
	//
	return y.X.Call("Q")[0].(rh.Q)
}

func (y *NoPanicYFID) Open(flag rh.Flag, intr rh.Intr) (err error) {
	//log.Println(y.String(), "Open")
	defer func() {
		if r := recover(); r != nil {
			err = rh.ErrIO // errors.NewPanic(r)
		}
	}()
	//
	return rh.Naturalize(errors.Unpack(y.X.Call("Open", flag, newXIntr(intr))[0]))
}

func (y *NoPanicYFID) Create(name string, flag rh.Flag, mode rh.Mode, perm rh.Perm) (fid rh.FID, err error) {
	//log.Println(y.String(), "Create")
	defer func() {
		if r := recover(); r != nil {
			fid, err = nil, rh.ErrIO // errors.NewPanic(r)
		}
	}()
	//
	r := y.X.Call("Create", name, flag, mode, perm)
	if err := errors.Unpack(r[1]); err != nil {
		return nil, rh.Naturalize(err)
	}
	return &NoPanicYFID{X: r[0].(circuit.X)}, nil
}

func (y *NoPanicYFID) Clunk() (err error) {
	//log.Println(y.String(), "Clunk")
	defer func() {
		if r := recover(); r != nil {
			err = rh.ErrIO
		}
	}()
	//
	r := y.X.Call("Clunk")
	err = errors.Unpack(r[0])
	if err != nil {
		return rh.Naturalize(err)
	}
	return nil
}

func (y *NoPanicYFID) Stat() (dir *rh.Dir, err error) {
	//log.Println(y.String(), "Stat")
	defer func() {
		if r := recover(); r != nil {
			dir, err = nil, rh.ErrIO // errors.NewPanic(r)
		}
	}()
	//
	r := y.X.Call("Stat")
	if err := errors.Unpack(r[1]); err != nil {
		return nil, rh.Naturalize(err)
	}
	return r[0].(*rh.Dir), nil
}

func (y *NoPanicYFID) Wstat(wdir *rh.Wdir) (err error) {
	//log.Println(y.String(), "Wstat")
	defer func() {
		if r := recover(); r != nil {
			err = rh.ErrIO // errors.NewPanic(r)
		}
	}()
	//
	return rh.Naturalize(errors.Unpack(y.X.Call("Wstat", wdir)[0]))
}

func (y *NoPanicYFID) Walk(wname []string) (fid rh.FID, err error) {
	//log.Println(y.String(), "Walk")
	defer func() {
		if r := recover(); r != nil {
			fid, err = nil, rh.ErrIO // errors.NewPanic(r)
		}
	}()
	//
	r := y.X.Call("Walk", wname)
	if err = errors.Unpack(r[1]); err != nil {
		return nil, rh.Naturalize(err)
	}
	return &NoPanicYFID{X: r[0].(circuit.X)}, nil
}

func (y *NoPanicYFID) Read(offset int64, count int, intr rh.Intr) (chunk rh.Chunk, err error) {
	//log.Println(y.String(), "Read")
	defer func() {
		if r := recover(); r != nil {
			chunk, err = nil, rh.ErrIO // errors.NewPanic(r)
		}
	}()
	//
	r := y.X.Call("Read", offset, count, newXIntr(intr))
	if err = errors.Unpack(r[1]); err != nil {
		return nil, rh.Naturalize(err)
	}
	return r[0].(rh.Chunk), nil
}

func (y *NoPanicYFID) Write(offset int64, data rh.Chunk, intr rh.Intr) (n int, err error) {
	//log.Println(y.String(), "Write")
	defer func() {
		if r := recover(); r != nil {
			n, err = 0, rh.ErrIO // errors.NewPanic(r)
		}
	}()
	//
	r := y.X.Call("Write", offset, data, newXIntr(intr))
	return r[0].(int), rh.Naturalize(errors.Unpack(r[1]))
}

func (y *NoPanicYFID) Remove() (err error) {
	//log.Println(y.String(), "Remove")
	defer func() {
		if r := recover(); r != nil {
			err = rh.ErrIO // errors.NewPanic(r)
		}
	}()
	//
	return rh.Naturalize(errors.Unpack(y.X.Call("Remove")[0]))
}

func (y *NoPanicYFID) Move(dir rh.FID, name string) (err error) {
	//log.Println(y.String(), "Move")
	defer func() {
		if r := recover(); r != nil {
			err = rh.ErrIO // errors.NewPanic(r)
		}
	}()
	//
	return rh.Naturalize(errors.Unpack(y.X.Call("Move", dir, name)[0]))
}

// The presence of the FID method ensures that NoPanicYFID conforms to FID at compile timerrors.
func (y *NoPanicYFID) FID() rh.FID {
	return y
}

// NoPanicYIntr
type NoPanicYIntr struct {
	rh.Intr
	sync.Mutex
	intr chan<- rh.Prompt
}

func NewNoPanicYIntr(xintr circuit.X) *NoPanicYIntr {
	intr := make(chan rh.Prompt, 1)
	y := &NoPanicYIntr{
		Intr: rh.Intr(intr),
		intr: intr,
	}
	go y.wait(xintr)
	return y
}

func (y *NoPanicYIntr) wait(xintr circuit.X) {
	defer func() {
		if r := recover(); r != nil {
			y.Clunk()
		}
	}()
	p := xintr.Call("Wait")[0].(rh.Prompt)
	//
	y.Lock()
	defer y.Unlock()
	//
	if y.intr == nil {
		return
	}
	y.intr <- p
	close(y.intr)
	y.intr = nil
}

func (y *NoPanicYIntr) Clunk() {
	y.Lock()
	defer y.Unlock()
	//
	if y.intr == nil {
		return
	}
	close(y.intr)
	y.intr = nil
}
