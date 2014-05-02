// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"encoding/json"
	"syscall"
)

type Stat struct {
	Cmd Cmd `json:"cmd"`
	Exit error `json:"exit"`
	State string `json:"state"`
}

func (s Stat) String() string {
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(0)
	}
	return string(b)
}

type RunState int

const (
	Unknown RunState = iota
	None
	Running
	Exited
	Stopped
	Signaled
	Continued
)

func (rs RunState) String() string {
	switch rs {
	case None:
		return "not executed yet"
	case Running:
		return "running"
	case Exited:
		return "exited"
	case Stopped:
		return "stopped"
	case Signaled:
		return "signaled"
	case Continued:
		return "continued"
	}
	return "unknown"
}

var (
	sigMap = map[string]syscall.Signal{
		"ABRT":   syscall.SIGABRT,
		"ALRM":   syscall.SIGALRM,
		"BUS":    syscall.SIGBUS,
		"CHLD":   syscall.SIGCHLD,
		"CONT":   syscall.SIGCONT,
		"FPE":    syscall.SIGFPE,
		"HUP":    syscall.SIGHUP,
		"ILL":    syscall.SIGILL,
		"INT":    syscall.SIGINT,
		"IO":     syscall.SIGIO,
		"IOT":    syscall.SIGIOT,
		"KILL":   syscall.SIGKILL,
		"PIPE":   syscall.SIGPIPE,
		"PROF":   syscall.SIGPROF,
		"QUIT":   syscall.SIGQUIT,
		"SEGV":   syscall.SIGSEGV,
		"STOP":   syscall.SIGSTOP,
		"SYS":    syscall.SIGSYS,
		"TERM":   syscall.SIGTERM,
		"TRAP":   syscall.SIGTRAP,
		"TSTP":   syscall.SIGTSTP,
		"TTIN":   syscall.SIGTTIN,
		"TTOU":   syscall.SIGTTOU,
		"URG":    syscall.SIGURG,
		"USR1":   syscall.SIGUSR1,
		"USR2":   syscall.SIGUSR2,
		"VTALRM": syscall.SIGVTALRM,
		"WINCH":  syscall.SIGWINCH,
		"XCPU":   syscall.SIGXCPU,
		"XFSZ":   syscall.SIGXFSZ,
	}
)
