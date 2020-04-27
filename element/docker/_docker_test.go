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

	ds "github.com/hoijui/circuit/client/docker"
)

func TestDocker(t *testing.T) {
	if _, err := Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	run := ds.Run{
		Image: "ubuntu",
		Path: "/bin/ls",
		Args: []string{"/"},
	}
	con, err := MakeContainer(run)
	if err != nil {
		t.Fatalf("make: %v", err)
	}
	con.Stdin().Close()
	peek, err := con.Wait()
	if err != nil {
		t.Fatalf("wait: %v", err)
	}
	fmt.Printf("%v\n", peek)
}
