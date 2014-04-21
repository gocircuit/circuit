// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package sel

import (
	"fmt"
	"bytes"
	"os"
	"os/exec"
	"path"

	"github.com/gocircuit/circuit/kit/fs/rh"
)

type result struct {
	Branch int
	Name   string
	Error  error
}

func waitOpenFile(branch int, name string, intr rh.Intr, report chan<- *result) {
	var r = &result{
		Branch: branch,
		Name:   name,
		Error:  OpenFile(name, intr),
	}
	//println(fmt.Sprintf("wait opened clause=%v name=%v error=%v", r.Branch, r.Name, r.Error))
	report <- r
}

// OpenFile tries to open the named local file for reading.
// On success, the returned error is nil.
// The errors os.ErrPermission or os.ErrNotExist, 
// if returned by the POSIX file open operation, are reported as is.
// All other non-nil errors represent other reasons for failing to open.
func OpenFile(file string, intr rh.Intr) error {
	wait := make(chan error, 1)
	cmd := exec.Command(getCircuitBinary(), "-sysopen", path.Clean(file))
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	go func() {
		wait <- cmd.Wait()
	}()
	//
	select {
	case err := <-wait:
		if stdout.String() == "not exist" {
			return os.ErrNotExist
		}
		if stdout.String() == "permission" {
			return os.ErrPermission
		}
		if stdout.String() != "" {
			return fmt.Errorf("open returned error: %s", stdout.String())
		}
		return err
	case <-intr:
		cmd.Process.Kill()
		return rh.ErrIntr
	}
}

func getCircuitBinary() string {
	if os.Getenv("CIRCUIT_BIN") != "" {
		return os.Getenv("CIRCUIT_BIN")
	}
	return os.Args[0]
}
