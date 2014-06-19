// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"os/exec"
)

func Init() (err error) {
	dkr, err = exec.LookPath("docker")
	if err != nil {
		return err
	}
	return nil
}

var dkr string

const StdBufferLen = 32e3
