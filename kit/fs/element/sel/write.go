// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package sel

import (
	"bufio"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/interruptible"
)

func getCircuitBinary() string {
	if os.Getenv("CIRCUIT_BIN") != "" {
		return os.Getenv("CIRCUIT_BIN")
	}
	return os.Args[0]
}

// FileWriter is an interruptible.Writer
type FileWriter struct {
	cmd   *exec.Cmd
	iw    interruptible.Writer
	kill  sync.Once
	Clunk func() // Clunk will be invoked if and when the FileWriter is closed.
}

func OpenFileWriter(name string, intr rh.Intr) (w *FileWriter, err error) {
	w = &FileWriter{
		cmd:   exec.Command(getCircuitBinary(), "-syswrite", path.Clean(name)),
	}
	// stderr is a control back-channel
	stderr, err := w.cmd.StderrPipe()
	if err != nil {
		panic(err) // system error
	}
	// stdin is the writing channel
	w.cmd.Stdin, w.iw = interruptible.BufferPipe(bufferCap)
	//
	if err = w.cmd.Start(); err != nil {
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
		w.cmd.Wait()
	}()
	select {
	case <-intr:
		w.cmd.Process.Kill()
		return nil, rh.ErrIntr
	case err := <-waitopen:
		if err != nil {
			w.cmd.Process.Kill()
			return nil, err
		}
		// Kill helper process on garbage collection
		runtime.SetFinalizer(w, 
			func(w2 *FileWriter) { 
				w2.Close()
			},
		)
		return w, nil
	}
}

func (w *FileWriter) WriteIntr(p []byte, intr rh.Intr) (_ int, err error) {
	defer func() {
		if err != nil {
			w.Close()
		}
	}()
	return w.iw.WriteIntr(p, intr)
}

func (w *FileWriter) Write(p []byte) (_ int, err error) {
	defer func() {
		if err != nil {
			w.Close()
		}
	}()
	return w.iw.Write(p)
}

func (w *FileWriter) Close() error {
	defer func() {
		defer func() {
			recover()
		}()
		w.kill.Do(func() {
			w.cmd.Process.Kill()
			if w.Clunk != nil {
				w.Clunk()
				w.Clunk = nil
			}
		})
	}()
	return w.iw.Close()
}

// DelayedWriteFile is an RH file backed by a FileWriter.
type DelayedWriteFile struct {
	w *FileWriter
}

func NewDelayedWriteFile(w *FileWriter) file.File {
	return &DelayedWriteFile{w: w}
}

func (f *DelayedWriteFile) Perm() rh.Perm {
	return 0222 // -w--w--w-
}

func (f *DelayedWriteFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	if flag.Attr != rh.WriteOnly && !flag.Truncate {
		return nil, rh.ErrPerm
	}
	return file.NewOpenInterruptibleWriterFile(f.w), nil
}

func (f *DelayedWriteFile) Remove() error {
	return rh.ErrPerm
}
