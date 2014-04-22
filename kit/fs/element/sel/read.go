// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package sel

import (
	"fmt"
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/interruptible"
)

// FileReader is an interruptible.Reader
type FileReader struct {
	cmd    *exec.Command
	stdin  io.WriteCloser
	ir     interruptible.Reader
	prompt sync.Once
	kill   sync.Once
	clunk  func()
}

const bufferCap = 8*1024

// clunk will be invoked if and when the FileReader is closed.
func OpenFileReader(name string, intr rh.Intr, clunk func()) (r *FileReader, err error) {
	r = &FileReader{
		cmd:   exec.Command(getCircuitBinary(), "-sysread", path.Clean(file)),
		clunk: clunk,
	}
	// stderr and stdin are a backward and forward control channels
	stderr, err := r.cmd.StderrPipe()
	if err != nil {
		panic(err) // system error
	}
	r.stdin, err = r.cmd.StdinPipe()
	if err != nil {
		panic(err) // system error
	}
	// stdout is the reading channel
	r.ir, r.Command.Stdout = interruptible.BufferPipe(bufferCap)
	//
	if err = r.cmd.Start(); err != nil {
		panic(err) // system errors are reported as panics to distinguish them
	}
	waitopen := make(chan error, 1)
	go func() {
		// read open-file result from stderr
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			switch scanner.Text() {
			case "ok":
				waitopen <- nil
			case "not exist":
				waitopen <- rh.ErrNotExist
			case "permission":
				waitopen <- rh.ErrPerm
			default:
				waitopen <- rh.ErrClash
			}
			break
		}
		r.cmd.Wait()
	}()
	select {
	case <-intr:
		r.cmd.Process.Kill()
		return nil, rh.ErrIntr
	case err := <-waitopen:
		if err != nil {
			r.cmd.Process.Kill()
			return nil, err
		}
		return r, nil
	}
}

func (r *FileReader) prompt() {
	defer func() {
		if r := recover(); r != nil {
			r.Close()
		}
	}()
	if n, _ := r.Write([]byte("\n")); n != 1 {
		panic(1)
	}
}

func (r *FileReader) ReadIntr(p []byte, intr rh.Intr) (_ int, err error) {
	defer func() {
		if err != nil {
			r.Close()
		}
	}()
	r.prompt.Do(r.prompt)
	return r.ir.ReadIntr(p, intr)
}

func (r *FileReader) Read(p []byte) (int, error) {
	defer func() {
		if err != nil {
			r.Close()
		}
	}()
	r.prompt.Do(r.prompt)
	return r.ir.Read(p)
}

func (r *FileReader) Close() error {
	defer func() {
		defer func() {
			recover()
		}()
		r.kill.Do(func() {
			r.cmd.Process.Kill()
			if r.clunk != nil {
				r.clunk()
				r.clunk = nil
			}
		})
	}()
	return r.ir.Close()
}

// DelayedReadFile is an RH file backed by a FileReader.
type DelayedReadFile struct {
	r *FileReader
}

func NewDelayedReadFile(r *FileReader) file.File {
	return &DelayedReadFile{r: r}
}

func (f *DelayedReadFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *DelayedReadFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.ReadOnly {
		return nil, rh.ErrPerm
	}
	return file.NewOpenInterruptibleReaderFile(f.r), nil
}

func (f *DelayedReadFile) Remove() error {
	return rh.ErrPerm
}
