// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"io"
	"os"
)

type Client struct {
}

func NewClient(mount string) (*Client, error) {
	x
}

func (c *Client) MakeChan(where string) Chan {
}

func (c *Client) MakeProc(where string) Proc {
}

func (c *Client) MakeSelect(where string) Select {
}
