// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package debug implements debugging utilities
package debug

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

// InstallTimeout panics the current process in ns time
func InstallTimeoutPanic(ns int64) {
	go func() {
		k := int(ns / 1e9)
		for i := 0; i < k; i++ {
			time.Sleep(time.Second)
			fmt.Fprintf(os.Stderr, "•%d/%d•\n", i, k)
		}
		//time.Sleep(time.Duration(ns))
		panic("process timeout")
	}()
}

func OnSignal(dfr func(os.Signal)) {
	go func() {
		//defer SavePanicTrace()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill)
		for sig := range ch {
			dfr(sig)
			os.Exit(1)
		}
	}()
}

// InstallCtrlCPanic installs a Ctrl-C signal handler that panics
func InstallCtrlCPanic() {
	go func() {
		//defer SavePanicTrace()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		for _ = range ch {
			panic("ctrlc")
			prof := pprof.Lookup("goroutine")
			prof.WriteTo(os.Stderr, 2)
			os.Exit(1)
		}
	}()
}

// InstallKillPanic installs a kill signal handler that panics
// From the command-line, this signal is agitated with kill -ABRT
func InstallKillPanic() {
	go func() {
		//defer SavePanicTrace()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Kill)
		for _ = range ch {
			prof := pprof.Lookup("goroutine")
			prof.WriteTo(os.Stderr, 2)
			os.Exit(1)
			//panic("sigkill")
		}
	}()
}

func SavePanicTrace() {
	r := recover()
	if r == nil {
		return
	}
	// Redirect stderr
	file, err := os.Create("panic")
	if err != nil {
		panic("dumper (no file) " + r.(fmt.Stringer).String())
	}
	syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	// TRY: defer func() { file.Close() }()
	panic("dumper " + r.(string))
}
