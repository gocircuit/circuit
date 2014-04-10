// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package xy

import (
	"runtime"
	"sync"
	"time"

	"github.com/gocircuit/circuit/sys/anchor"
	use "github.com/gocircuit/circuit/use/anchorfs"
	"github.com/gocircuit/circuit/use/circuit"
)

func init() {
	circuit.RegisterValue(&XSys{})
	circuit.RegisterValue(&XDir{})
	circuit.RegisterValue(&XKeepalive{})
}

// XKeepalive
type XKeepalive struct {
	sync.Once
	f func()
}

func NewXKeepalive(f func()) *XKeepalive {
	xkeep := &XKeepalive{f: f}
	runtime.SetFinalizer(xkeep, func(x *XKeepalive) {
		x.Once.Do(x.f)
	})
	return xkeep
}

// XSys
type XSys anchor.System

func (s *XSys) Remove(pay anchor.Payload) {
	(*anchor.System)(s).Remove(pay)
}

func (s *XSys) Create(dirs []string, pay anchor.Payload) (xkeep circuit.X, err error) {
	if err = (*anchor.System)(s).Create(dirs, pay); err != nil {
		return nil, err
	}
	return circuit.Ref(NewXKeepalive(func() {
		s.Remove(pay)
	})), nil
}

func (s *XSys) OpenFile(fullpath string) (*anchor.File, error) {
	return (*anchor.System)(s).OpenFile(fullpath)
}

func (s *XSys) OpenDir(fullpath string) (circuit.X, error) {
	d, err := (*anchor.System)(s).OpenDir(fullpath)
	return circuit.Ref((*XDir)(d)), err
}

// XDir
type XDir anchor.Dir

func (d *XDir) Path() string {
	return (*anchor.Dir)(d).Path()
}

func (d *XDir) List() (rev use.Rev, files, dirs []string) {
	return (*anchor.Dir)(d).List()
}

func (d *XDir) Change(sinceRev use.Rev) (rev use.Rev, files, dirs []string) {
	return (*anchor.Dir)(d).Change(sinceRev)
}

func (d *XDir) ChangeExpire(sinceRev use.Rev, expire time.Duration) (rev use.Rev, files, dirs []string, err error) {
	return (*anchor.Dir)(d).ChangeExpire(sinceRev, expire)
}

func (d *XDir) OpenDir(name string) (circuit.X, error) {
	d_, err := (*anchor.Dir)(d).OpenDir(name)
	return circuit.Ref((*XDir)(d_)), err
}

func (d *XDir) OpenFile(name string) (*anchor.File, error) {
	return (*anchor.Dir)(d).OpenFile(name)
}
