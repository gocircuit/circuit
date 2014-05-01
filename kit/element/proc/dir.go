// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"fmt"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// ProcDir 
type ProcDir struct {
	name string
	rh.FID
	dir *dir.Dir
	p   *Proc
	rmv struct {
		sync.Mutex
		rmv func() // remove from parent
	}
}

func NewDir(name string, rmv func()) *ProcDir {
	d := &ProcDir{
		dir: dir.NewDir(),
		p:   MakeProc(),
	}
	d.rmv.rmv = rmv
	d.FID = d.dir.FID()
	d.dir.AddChild("help", file.NewFileFID(file.NewByteReaderFile(
		func() []byte {
			return []byte(d.Help())
		}),
	))
	d.dir.AddChild("env", file.NewFileFID(NewEnvFile(d.p)))
	//
	d.dir.AddChild("stdin", file.NewFileFID(NewStdinFile(d.p)))
	d.dir.AddChild("stdout", file.NewFileFID(NewStdoutFile(d.p)))
	d.dir.AddChild("stderr", file.NewFileFID(NewStderrFile(d.p)))
	//
	d.dir.AddChild("cmd", file.NewFileFID(NewCmdFile(d.p)))
	d.dir.AddChild("wait", file.NewFileFID(NewWaitFile(d.p)))
	d.dir.AddChild("trywait", file.NewFileFID(NewTryWaitFile(d.p)))
	d.dir.AddChild("signal", file.NewFileFID(NewSignalFile(d.p)))
	//
	d.dir.AddChild("error", file.NewFileFID(d.p.ErrorFile))

	return d
}

func (s *ProcDir) Walk(wname []string) (rh.FID, error) {
	if len(wname) > 0 {
		return s.FID.Walk(wname)
	}
	return s, nil
}

func (d *ProcDir) Remove() error {
	d.rmv.Lock()
	defer d.rmv.Unlock()
	if err := d.p.ClunkIfNotBusy(); err != nil {
		return rh.ErrBusy
	}
	if d.rmv.rmv != nil {
		d.rmv.rmv()
	}
	d.rmv.rmv = nil
	return nil
}

func (d *ProcDir) Help() string {
	return fmt.Sprintf(dirHelpFormat, d.name)
}

const dirHelpFormat = `
	This is the control directory for a circuit process named: %s

	A circuit process is like a handle for a unique OS process execution.

RUN

	To start a new process, write a JSON-encoded execution description to "run", 
	like so

		echo << EOF > run
		{
			"env": ["PATH=/usr/bin", "KEY=value"],
			"path": "/bin/ls",
			"args": ["-a", "/"]
		}
		EOF

	Only "path" is mandatory. Writing to "run" will not block: starting
	a process is an asynchronous operation. To be notified of the exit
	of a process, see WAITING below.

STDIN, STDOUT, STDERR

	The files "stdin", "stdout" and "stderr" are directly connected
	to the standard input, output and error file descriptors of the
	executing OS process.

		echo ls > stdin
		cat stdout
		cat stderr

	They can be opened for writing, reading and reading, respectively,
	at any time before or after the execution of the process (see RUN).

	Whenever data is not available for reading on standard output or 
	error, or the process is not accepting new writes to standard input,
	respective operations on the files will block.

	The user is responsible for emptying the contents of "stdout"
	and "stderr". A small buffer is included for convenience, but if
	too much output is left unread, the underlying process will block.

WAITING

	To wait for the exit of an OS processes, executed with RUN, open
	and read from "waitexit".

		cat waitexit

	Opening "waitexit" will block until the process is still running.
	When the process exits, open will unblock and the exit conditions
	will be readable from "waitexit". On success, the contents of
	"waitexit" will be empty. Otherwise, an error message will be 
	readable.

	If the process is already dead, open will succeed right away.
	Multiple and concurrent invocations of "waitexit" are allowed.

SIGNAL

	To send an OS signal to a running process asynchronously,
	write the POSIX name of the signal or the textual representation 
	of its number to "signal".

		echo KILL > signal
		echo 9 > signal

STAT

	Reading from "stat" will return the current execution state of
	the process. One of:
		not executed yet
		running
		exited
		stopped
		signaled
		continued
		unknown

ENV

	Reading "env" asynchronously returns the OS environment at the
	hosting circuit.

		cat env

ERRORS

	Unsuccessful operations with the special in this directory will
	return file system errors. These errors are standardized in POSIX
	and are not descriptive for our purposes. To remedy this, after
	every file manipulation that returns a file system error, a detailed
	error message will be readable from the "error" file until the next
	file manipulation.

		cat error

`
