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

// client -> circuit -> weave=knot -> chan, proc

// CircuitID opaquely identifies a circuit worker
type CircuitID string

// Client is a live session with the circuit mount point on the local machine.
type Client struct {
	mount *Dir
}

// NewClient creates a new client for the circuit environment.  It opens the
// circuit mount point and keeps it open (preventing it from being unmounted)
// for the life of the returned Client object.
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

// Circuits asynchronously returns a list of known live circuits.
func (c *Client) Circuits() ([]CircuitID, error) {
	d, err := os.Open(c.mount.Path())
	if err != nil {
		return nil, err
	}
	defer d.Close()
	children, err := d.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	var w []CircuitID
	for _, name := range children {
		if name == "" || name[0] != 'X' {
			continue
		}
		w = append(w, CircuitID(name))
	}
	return w, nil
}

// Circuit returns a controller for the circuit instance with the given id.
func (c *Client) Circuit(id CircuitID) (*Circuit, error) {
	return openCircuit(c, id)
}
