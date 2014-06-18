// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"sync"

	dkr "github.com/fsouza/go-dockerclient"
)

func Connect(endpoint string) (err error) {
	cli.Lock()
	defer cli.Unlock()
	cli.Client, err = dkr.NewClient(endpoint)
	return
}

var (
	cli struct {
		sync.Mutex
		*dkr.Client
	}
)

func client() *dkr.Client {
	cli.Lock()
	defer cli.Unlock()
	return cli.Client
}
