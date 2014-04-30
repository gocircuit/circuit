// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type Command struct {
	Env  []string `json:"env"`
	Path string   `json:"path"`
	Args []string `json:"args"`
}

const procDir = "Proc"

// ErrExit represents any non-zero process exit code.
var ErrExit = errors.New("circuit process exit error")

type Proc struct {
	local string
	dir *Dir
}

func openProc(local string) *Proc {
	p = &Proc{local: local}
	if p.dir, err = OpenDir(p.Path()); err != nil {
		panic(err)
	}
	return p
}

// Path returns the path of this Process element in the local circuit file system.
func (p *Proc) Path() string {
	return p.local
}

// Start …
func (p *Proc) Start(cmd Command) error {
	b, err := json.Marshal(cmd)
	if err != nil {
		panic(0)
	}
	return ioutil.WriteFile(path.Join(p.Path(), "start"), b, 0222)
}

// Wait …
func (p *Proc) Wait() error {
	b, err := ioutil.ReadFile(path.Join(p.Path(), "waitexit"))
	if os.IsNotExist(err) { // a missing file indicates a dead circuit worker; we panic for those by convention
		panic(err)
	}
	if err != nil { // other errors are Process element specific; we report them traditionally
		return err
	}
	if len(b) == 0 {
		return nil
	}
	return ErrExit
}
