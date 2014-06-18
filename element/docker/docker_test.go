// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"fmt"
	"testing"

	//dkr "github.com/fsouza/go-dockerclient"
)

func TestDocker(t *testing.T) {
	if err := Connect("unix:///var/run/docker.sock"); err != nil {
		t.Fatalf("connect: %v", err)
	}
	run := Run{
		Image: "b6b9590f1a97",
		Path: "/bin/ls",
		Args: []string{"/"},
	}
	con, err := MakeContainer(run)
	if err != nil {
		t.Fatalf("make: %v", err)
	}
	fmt.Printf("%v\n", con)
}
