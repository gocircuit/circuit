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
	id string
	dir *dir
	anchor *Anchor
}

func openCircuit(client *Client, id string) (c *Circuit) {
	c = &Circuit{
		client: client,
		id: id,
	}
	var err error
	if c.dir, err = openDir(c.Path()); err != nil {
		panic(err)
	}
	if c.anchor, err = openAnchor(path.Join(c.Path(), anchorDir)); err != nil {
		panic(err)
	}
	return c
}

// Path returns the path of this circuit in the local file system.
func (c *Circuit) Path() string {
	return path.Join(c.client.Path(), string(c.id))
}

const anchorDir = "anchor"

// Anchor
func (c *Circuit) Anchor(walk ...string) *Anchor {
	return c.anchor.Anchor(walk...)
}
