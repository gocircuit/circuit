// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package file

import (
	"io"

	"github.com/gocircuit/circuit/kit/interruptible"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// NewOpenWriterFile returns an FID corresponding to an open reader file.
func NewOpenWriterFile(w io.WriteCloser) rh.FID {
	file := &OpenWriterFile{
		openFile: newOpenFile(),
	}
	ir, iw := interruptible.Pipe()
	go func() {
		io.Copy(w, ir)
		ir.Close()
		w.Close()
	}()
	file.w = iw
	return file
}

func NewOpenInterruptibleWriterFile(w interruptible.Writer) rh.FID {
	return &OpenWriterFile{
		openFile: newOpenFile(),
		w:        w,
	}
}

// NewOpenPipeWriterFile returns an opened reader file, connected to a pipe writer on the other, also returned.
func NewOpenPipeWriterFile() (file rh.FID, r interruptible.Reader) {
	w := &OpenWriterFile{
		openFile: newOpenFile(),
	}
	r, w.w = interruptible.Pipe()
	return w, r
}

// OpenWriterFile is an FID representing an open write-only file.
type OpenWriterFile struct {
	*openFile
	w interruptible.Writer
}

func (w *OpenWriterFile) Walk(wname []string) (rh.FID, error) {
	if len(wname) > 0 {
		return nil, rh.ErrClash
	}
	println("trying to clone open fid")
	return nil, rh.ErrClash
}

func (w *OpenWriterFile) Clunk() (err error) {
	if err = w.w.Close(); err != nil {
		return rh.ErrClash
	}
	w.openFile.Clunk()
	return
}

func (w *OpenWriterFile) Stat() (dir *rh.Dir, err error) {
	dir, _ = w.openFile.Stat()
	dir.Perm = 0222 // -w--w--w-
	return dir, nil
}

func (w *OpenWriterFile) Write(_ int64, data rh.Chunk, intr rh.Intr) (n int, err error) {
	if data == nil {
		return 0, nil
	}
	chnk, ok := data.(rh.ByteChunk)
	if !ok {
		return 0, rh.ErrClash
	}
	n, err = w.w.WriteIntr(chnk, intr)
	if err != nil {
		return n, rh.ErrIO
	}
	return n, nil
}
