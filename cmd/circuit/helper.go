// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"fmt"
	"io"
	"os"
)

func runAsMain(run func() error) {
	if err := run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func mainSysOpenRead(file string) {
	runAsMain(func() error {
		return sysOpenRead(file)
	})
}

func mainSysOpenWrite(file string) {
	runAsMain(func() error {
		return sysOpenWrite(file)
	})
}

func sysOpenRead(file string) error {
	f, err := os.OpenFile(file, os.O_RDONLY, 0444)
	if os.IsNotExist(err) {
		println("not exist")
		return err
	}
	if os.IsPermission(err) {
		println("permission")
		return err
	}
	if err != nil {
		println("unknown")
		return err
	}
	println("ok")
	os.Stderr.Sync()
	// It is imperative that we don't invoke f's Read operation  (which would
	// commit channel receivers) before receiving a prompt from the parent process.
	if _, err := fmt.Scanln(); err != nil {
		return err
	}
	// Start reading from file, buffering, and pumping into standard output back to the parent process.
	if _, err := io.Copy(os.Stdout, f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func sysOpenWrite(file string) error {
	f, err := os.OpenFile(file, os.O_WRONLY, 0222)
	if os.IsNotExist(err) {
		println("not exist")
		return err
	}
	if os.IsPermission(err) {
		println("permission")
		return err
	}
	if err != nil {
		println("unknown")
		return err
	}
	println("ok")
	os.Stderr.Sync()
	if _, err := io.Copy(f, os.Stdin); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
