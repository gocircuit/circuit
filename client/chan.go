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

const channelDir = "chan"

type channel struct {
	namespace *Namespace
	name      string
	dir       *Dir
}

func makeChan(namespace *Namespace, name string, cap_ int) (c *channel, err error) {
	c = &channel{
		namespace: namespace,
		name:      name,
	}
	if err = os.Mkdir(c.Path(), 0777); err != nil {
		return nil, err
	}
	if err = ioutil.WriteFile(path.Join(c.Path(), "cap"), []byte(strconv.Itoa(cap_)), 0222); err != nil {
		os.Remove(c.Path())
		return nil, err
	}
	if c.dir, err = OpenDir(c.Path()); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *channel) Path() string {
	return path.Join(c.namespace.Path(), channelDir, c.name)
}

func (c *channel) Send() io.WriteCloser {
	f, err := os.OpenFile(path.Join(c.Path(), "send"), os.O_WRONLY, 0222)
	if err != nil {
		panic(err)
	}
	return f
}

func (c *channel) Recv() io.ReadCloser {
	f, err := os.OpenFile(path.Join(c.Path(), "recv"), os.O_RDONLY, 0444)
	if err != nil {
		return nil
	}
	return f
}

func (c *channel) WaitSend() {
	f, err := os.OpenFile(path.Join(c.Path(), "waitsend"), os.O_RDONLY, 0444)
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

func (c *channel) WaitRecv() {
	f, err := os.OpenFile(path.Join(c.Path(), "waitrecv"), os.O_RDONLY, 0444)
	if err != nil {
		return nil
	}
	defer f.Close()
}

func (c *channel) TrySend() io.WriteCloser {
	f, err := os.OpenFile(path.Join(c.Path(), "trysend"), os.O_WRONLY, 0222)
	if err != nil {
		return nil
	}
	return f
}

func (c *channel) TryRecv() io.ReadCloser {
	f, err := os.OpenFile(path.Join(c.Path(), "tryrecv"), os.O_RDONLY, 0444)
	if err != nil {
		return nil
	}
	return f
}

func (c *channel) Close() {
	if err := ioutil.WriteFile(path.Join(c.Path(), "close"), []byte("close"), 0222); err != nil {
		panic(err)
	}
}
