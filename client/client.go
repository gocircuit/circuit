// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"os"
	"path"
)

// client -> worker -> namespace => chan,sync,proc

// WorkerID opaquely identifies a circuit worker
type WorkerID string

// Client is a live session with the circuit mount point on the local machine.
type Client struct {
	mount *Dir
}

// NewClient creates a new client for the circuit environment.  It opens the
// circuit mount point and keeps it open (preventing it from being unmounted)
// for the life of the return Client object.
func NewClient(mount string) (c *Client, err error) {
	c = &Client{}
	if c.mount, err = OpenDir(path.Clean(mount)); err != nil {
		return nil, err
	}
	return c, nil
}

// Path returns the local mount-point of the circuit file system used by this client.
func (c *Client) Path() string {
	return c.mount.Path()
}

// Workers asynchronously returns a list of known live workers.
func (c *Client) Workers() ([]WorkerID, error) {
	d, err := os.Open(c.mount.Path())
	if err != nil {
		return nil, err
	}
	defer d.Close()
	children, err := d.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	var w []WorkerID
	for _, name := range children {
		if name == "" || name[0] != 'X' {
			continue
		}
		w = append(w, WorkerID(name))
	}
	return w, nil
}

// Worker returns a controller for the requested circuit worker.
func (c *Client) Worker(worker WorkerID) (*Worker, error) {
	return newWorker(c, worker)
}
