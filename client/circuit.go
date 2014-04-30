// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"path"
)

// Circuit is a client for a specific circuit worker's distributed control facilities.
type Circuit struct {
	client *Client
	id CircuitID
	dir *Dir
	anchor *Anchor
}

func openCircuit(client *Client, id CircuitID) (c *Circuit, err error) {
	c = &Circuit{
		client: client,
		id: id,
	}
	if c.dir, err = OpenDir(c.Path()); err != nil {
		return nil, err
	}
	if c.anchor, err = openAnchor(path.Join(c.Path(), anchorDir)); err != nil {
		return nil, err
	}
	return c, nil
}

// Path returns the path of this circuit in the local file system.
func (c *Circuit) Path() string {
	return path.Join(c.client.Path(), string(c.id))
}

const anchorDir = "anchor"

// UseAnchor
func (c *Circuit) UseAnchor(walk []string) (_ *Anchor, err error) {
	return c.anchor.Use(walk)
}
