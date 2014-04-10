// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package xy

import (
	"github.com/gocircuit/circuit/sys/anchor"
	use "github.com/gocircuit/circuit/use/anchorfs"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
	"time"
)

// YKeepalive
type YKeepalive circuit.X

// YSys
type YSys struct {
	circuit.X
}

func (s YSys) Remove(pay anchor.Payload) {
	s.X.Call("Remove", pay)
}

func (s YSys) Create(dirs []string, pay anchor.Payload) (ykeep YKeepalive, err error) {
	r := s.X.Call("Create", dirs, pay)
	if err = e.Unpack(r[1]); err != nil {
		return nil, err
	}
	return YKeepalive(r[0].(circuit.X)), nil
}

func (s YSys) OpenFile(fullpath string) (*anchor.File, error) {
	r := s.X.Call("OpenFile", fullpath)
	if err := e.Unpack(r[1]); err != nil {
		return nil, err
	}
	return r[0].(*anchor.File), nil
}

func (s YSys) OpenDir(fullpath string) (YDir, error) {
	r := s.X.Call("OpenDir", fullpath)
	if err := e.Unpack(r[1]); err != nil {
		return YDir{nil}, err
	}
	return YDir{r[0].(circuit.X)}, nil
}

// YDir
type YDir struct {
	circuit.X
}

func (d YDir) Path() string {
	return d.X.Call("Path")[0].(string)
}

func (d YDir) List() (rev use.Rev, files, dirs []string) {
	r := d.X.Call("List")
	return r[0].(use.Rev), r[1].([]string), r[2].([]string)
}

func (d YDir) Change(sinceRev use.Rev) (rev use.Rev, files, dirs []string) {
	r := d.X.Call("Change", sinceRev)
	return r[0].(use.Rev), r[1].([]string), r[2].([]string)
}

func (d YDir) ChangeExpire(sinceRev use.Rev, expire time.Duration) (rev use.Rev, files, dirs []string, err error) {
	r := d.X.Call("ChangeExpire", sinceRev, expire)
	if err := e.Unpack(r[3]); err != nil {
		return 0, nil, nil, err
	}
	return r[0].(use.Rev), r[1].([]string), r[2].([]string), nil
}

func (d YDir) OpenDir(name string) (YDir, error) {
	r := d.X.Call("OpenDir", name)
	if err := e.Unpack(r[1]); err != nil {
		return YDir{}, err
	}
	return YDir{r[0].(circuit.X)}, nil
}

func (d YDir) OpenFile(name string) (*anchor.File, error) {
	r := d.X.Call("OpenFile", name)
	if err := e.Unpack(r[1]); err != nil {
		return nil, err
	}
	return r[0].(*anchor.File), nil
}
