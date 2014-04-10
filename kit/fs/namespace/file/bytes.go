// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package file

import (
	"bytes"

	"github.com/gocircuit/circuit/kit/iomisc"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// ByteReaderFile is a read-only file, backed by a slice of bytes.
type ByteReaderFile struct {
	payload func() []byte
}

func NewByteReaderFile(payload func() []byte) File {
	return &ByteReaderFile{payload}
}

func NewStringReaderFile(payload func() string) File {
	return &ByteReaderFile{
		func() []byte {
			return []byte(payload())
		},
	}
}

func (f *ByteReaderFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *ByteReaderFile) Open(rh.Flag, rh.Intr) (rh.FID, error) {
	return NewOpenReaderFile(iomisc.ReaderNopCloser(bytes.NewBuffer(f.payload()))), nil
}

func (f *ByteReaderFile) Remove() error {
	return rh.ErrPerm
}
