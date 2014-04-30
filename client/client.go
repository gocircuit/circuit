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

// Client is a live session with the circuit mount point on the local machine.
type Client struct {
	mount *dir
}

// Attach creates a new client for the circuit environment.  It opens the
// circuit mount point and keeps it open (preventing it from being unmounted)
// for the life of the returned Client object.
func Attach(mount string) (c *Client) {
	c = &Client{}
	var err error
	if c.mount, err = openDir(path.Clean(mount)); err != nil {
		panic(err)
	}
	return c
}

// Path returns the local mount-point of the circuit file system used by this client.
func (c *Client) Path() string {
	return c.mount.Path()
}

// Circuits asynchronously returns a list of known live circuits.
func (c *Client) Circuits() []string {
	d, err := os.Open(c.mount.Path())
	if err != nil {
		panic(err)
	}
	defer d.Close()
	children, err := d.Readdirnames(0)
	if err != nil {
		panic(err)
	}
	var w []string
	for _, name := range children {
		if name == "" || name[0] != 'X' {
			continue
		}
		w = append(w, name)
	}
	return w
}

// Circuit returns a controller for the circuit instance with the given id.
func (c *Client) Circuit(id string) *Circuit {
	return openCircuit(c, id)
}
