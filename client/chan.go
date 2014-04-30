// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

const chanDir = "chan"

// Chan is a handle for a circuit channel.
type Chan struct {
	local string
	dir *Dir
}

func openChan(local string) (c *Chan) {
	c = &Chan{local: local}
	var err error
	if c.dir, err = OpenDir(c.Path()); err != nil {
		panic(err)
	}
	return c
}

// Path returns the path of this Chan in the local circuit file system.
func (c *Chan) Path() string {
	return c.local
}

// SetCap must be invoked once on a Chan before the Chan can be used.
func (c *Chan) SetCap(n int) error {
	return ioutil.WriteFile(path.Join(c.Path(), "cap"), []byte(strconv.Itoa(n)), 0222)
}

// Send …
func (c *Chan) Send() io.WriteCloser {
	f, err := os.OpenFile(path.Join(c.Path(), "send"), os.O_WRONLY, 0222)
	if err != nil {
		panic(err)
	}
	return f
}

// Recv …
func (c *Chan) Recv() io.ReadCloser {
	f, err := os.OpenFile(path.Join(c.Path(), "recv"), os.O_RDONLY, 0444)
	if err != nil {
		return nil
	}
	return f
}

// Close closes the circuit Chan.
func (c *Chan) Close() error {
	return ioutil.WriteFile(path.Join(c.Path(), "close"), []byte(`"close"`), 0222)
}
