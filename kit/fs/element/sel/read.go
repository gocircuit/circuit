// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package sel

import (
	"bufio"
	"io"
	"os/exec"
	"path"
	"runtime"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/interruptible"
)

// FileReader is an interruptible.Reader
type FileReader struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	ir     interruptible.Reader
	prompt sync.Once
	kill   sync.Once
	Clunk  func() // Non-nil Clunk will be invoked if and when the FileReader is closed.
}

const bufferCap = 8*1024

// OpenFileReader tries to open the named file for reading, potentially blocking for a while.
func OpenFileReader(name string, intr rh.Intr) (r *FileReader, err error) {
	r = &FileReader{
		cmd: exec.Command(getCircuitBinary(), "-sysread", path.Clean(name)),
	}
	// the helper process communicates control messages back to us on stderr
	stderr, err := r.cmd.StderrPipe()
	if err != nil {
		panic(err) // system error
	}
	r.stdin, err = r.cmd.StdinPipe()
	if err != nil {
		panic(err) // system error
	}
	// stdout is the reading channel
	r.ir, r.cmd.Stdout = interruptible.BufferPipe(bufferCap)
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
		// Kill helper process on garbage collection
		runtime.SetFinalizer(r, 
			func(r2 *FileReader) { 
				r2.Close()
			},
		)
		println("opened file reader->")
		return r, nil
	}
}

func (r *FileReader) commit() {
	r.prompt.Do(
		func() {
			defer func() {
				if p := recover(); p != nil {
					r.Close()
				}
			}()
			if n, _ := r.stdin.Write([]byte("\n")); n != 1 {
				panic(1)
			}
		},
	)
}

func (r *FileReader) ReadIntr(p []byte, intr rh.Intr) (_ int, err error) {
	defer func() {
		if err != nil {
			r.Close()
		}
	}()
	r.commit()
	return r.ir.ReadIntr(p, intr)
}

func (r *FileReader) Read(p []byte) (_ int, err error) {
	defer func() {
		if err != nil {
			r.Close()
		}
	}()
	r.commit()
	return r.ir.Read(p)
}

func (r *FileReader) Close() error {
	defer func() {
		defer func() {
			recover()
		}()
		r.kill.Do(func() {
			r.cmd.Process.Kill()
			if r.Clunk != nil {
				r.Clunk()
				r.Clunk = nil
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
