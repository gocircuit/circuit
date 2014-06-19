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

func Dial(endpoint string) (err error) {
	config.Lock()
	defer config.Unlock()
	if _, err = dkr.NewClient(endpoint); err != nil {
		return err
	}
	config.endpoint = endpoint
	return
}

var config struct {
	sync.Mutex
	endpoint string
}

func dial() (cli *dkr.Client, err error) {
	config.Lock()
	endpoint := config.endpoint
	config.Unlock()
	cli, err = dkr.NewClient(endpoint)
	return
}
