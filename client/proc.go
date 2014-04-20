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

const procDir = "proc"

var ErrExit = errors.New("circuit process exit error")

type proc struct {
	namespace *Namespace
	name      string
	dir       *Dir
}

func makeProc(namespace *Namespace, name string) (p *proc, err error) {
	p = &proc{
		namespace: namespace,
		name:      name,
	}
	if err = os.Mkdir(p.Path(), 0777); err != nil {
		return nil, err
	}
	if p.dir, err = OpenDir(p.Path()); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *proc) Path() string {
	return path.Join(p.namespace.Path(), procDir, p.name)
}

func (p *proc) Start(cmd Command) error {
	b, err := json.Marshal(cmd)
	if err != nil {
		panic(0)
	}
	return ioutil.WriteFile(path.Join(p.Path(), "start"), b, 0222)
}

func (p *proc) Wait() error {
	b, err := ioutil.ReadFile(path.Join(p.Path(), "waitexit"))
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	return ErrExit
}
