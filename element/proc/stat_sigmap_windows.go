// +build windows
// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import "syscall"

var (
	sigMap = map[string]syscall.Signal{
		"ABRT": syscall.SIGABRT,
		"ALRM": syscall.SIGALRM,
		"BUS":  syscall.SIGBUS,
		"FPE":  syscall.SIGFPE,
		"HUP":  syscall.SIGHUP,
		"ILL":  syscall.SIGILL,
		"INT":  syscall.SIGINT,
		"KILL": syscall.SIGKILL,
		"PIPE": syscall.SIGPIPE,
		"QUIT": syscall.SIGQUIT,
		"SEGV": syscall.SIGSEGV,
		"TERM": syscall.SIGTERM,
		"TRAP": syscall.SIGTRAP,
	}
)
