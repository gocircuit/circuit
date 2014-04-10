// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func mount(dir string) (fd int, serr string) {
	// Load OSXFUSE
	var err error
	if _, err = getvfsbyname("fusefs"); err != nil {
		const loader = "/Library/Filesystems/osxfusefs.fs/Support/load_osxfusefs"
		if _, err = os.Stat(loader); err != nil {
			return 0, "cannot find load_fusefs"
		}
		if err = exec.Command("/bin/sh", "-c", loader).Run(); err != nil {
			return 0, err.Error()
		}
		if _, err = getvfsbyname("osxfusefs"); err != nil {
			return 0, err.Error()
		}
	}

	// Look for available FUSE device
	var fusedev *os.File
	for i := 0; ; i++ {
		name := fmt.Sprintf("/dev/osxfuse%d", i)
		if _, err = os.Stat(name); err != nil {
			return 0, "no available fuse devices"
		}
		if fusedev, err = os.OpenFile(name, os.O_RDWR, 0); err == nil {
			break
		}
	}
	fd2, err := syscall.Dup(int(fusedev.Fd()))
	if err != nil {
		fusedev.Close()
		return 0, err.Error()
	}

	// Mount
	const (
		mounterDir  = "/Library/Filesystems/osxfusefs.fs/Support"
		mounterName = "mount_osxfusefs"
		mounter     = mounterDir + "/" + mounterName
	)

	cmd := exec.Command(mounter, "-o", "iosize=4096", strconv.Itoa(3), dir)
	cmd.Dir = mounterDir
	cmd.Env = []string{
		"MOUNT_FUSEFS_CALL_BY_LIB=",
		"MOUNT_FUSEFS_DAEMON_PATH=" + mounter,
	}
	cmd.ExtraFiles = []*os.File{fusedev}

	if err = cmd.Start(); err != nil {
		fusedev.Close()
		syscall.Close(fd2)
		return 0, fmt.Sprintf("exec mount_osxfusefs (%s)", err)
	}
	go func() {
		defer func() {
			recover()
		}()
		cmd.Wait()
	}()

	return fd2, ""
}
