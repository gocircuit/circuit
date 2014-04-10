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

// NewOpenReaderFile returns an FID corresponding to an open reader file.
func NewOpenReaderFile(r io.ReadCloser) rh.FID {
	file := &OpenReaderFile{
		openFile: newOpenFile(),
	}
	ir, iw := interruptible.Pipe()
	go func() {
		io.Copy(iw, r)
		iw.Close()
		r.Close()
	}()
	file.r = ir
	return file
}

// NewOpenInterruptibleReaderFile
func NewOpenInterruptibleReaderFile(r interruptible.Reader) rh.FID {
	return &OpenReaderFile{
		openFile: newOpenFile(),
		r:        r,
	}
}

// NewOpenPipeReaderFile returns an opened reader file, connected to a pipe writer on the other, also returned.
func NewOpenPipeReaderFile() (file rh.FID, w interruptible.Writer) {
	r := &OpenReaderFile{
		openFile: newOpenFile(),
	}
	r.r, w = interruptible.Pipe()
	return r, w
}

// ··················································································

// OpenReaderFile is an FID representing an open read-only file.
type OpenReaderFile struct {
	*openFile
	r interruptible.Reader
}

func (r *OpenReaderFile) Walk(wname []string) (rh.FID, error) {
	if len(wname) > 0 {
		return nil, rh.ErrClash
	}
	println("trying to clone open fid")
	return nil, rh.ErrClash
}

func (r *OpenReaderFile) Clunk() (err error) {
	if err = r.r.Close(); err != nil {
		err = rh.ErrClash
	}
	r.openFile.Clunk()
	return
}

func (r *OpenReaderFile) Stat() (dir *rh.Dir, err error) {
	dir, _ = r.openFile.Stat()
	dir.Perm = 0444 // r--r--r--
	return dir, nil
}

func (r *OpenReaderFile) Read(_ int64, count int, intr rh.Intr) (chunk rh.Chunk, err error) {
	var buf = make(rh.ByteChunk, count)
	n, err := r.r.ReadIntr(buf, intr)
	if err != nil {
		if err == io.EOF {
			err = nil // Don't report EOF, since FUSE doesn't like it on regular files.
		} else {
			err = rh.ErrIO
		}
		return buf[:n], err
	}
	return buf[:n], nil
}
