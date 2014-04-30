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
	dir *dir
}

func openProc(local string) (p *Proc) {
	p = &Proc{local: local}
	var err error
	if p.dir, err = openDir(p.Path()); err != nil {
		panic(err)
	}
	return p
}

// Path returns the path of this Process element in the local circuit file system.
func (p *Proc) Path() string {
	return p.local
}

// Start …
func (p *Proc) Start(cmd Command) {
	b, err := json.Marshal(cmd)
	if err != nil {
		panic(0)
	}
	if err = ioutil.WriteFile(path.Join(p.Path(), "cmd"), b, 0222); err != nil {
		panic(err)
	}
}

type stat struct {
	Exit string `json:"exit"`
	State string `json:"state"`
}

// Wait …
func (p *Proc) Wait() error {
	b, err := ioutil.ReadFile(path.Join(p.Path(), "wait"))
	if os.IsNotExist(err) { // a missing file indicates a dead circuit worker; we panic for those by convention
		panic(err)
	}
	if err != nil { // other errors are element specific; we report them traditionally
		return err
	}
	var stat stat
	if err = json.Unmarshal(b, &stat); err != nil {
		panic(err)
	}
	if stat.Exit != "" {
		return ErrExit
	}
	return nil
}
