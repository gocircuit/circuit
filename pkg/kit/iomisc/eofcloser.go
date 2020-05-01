// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package iomisc

import (
	"io"
)

// ReaderEOF returns an io.Reader which will return an io.EOF error as
// soon as the reader r is empty.
func ReaderEOF(r io.Reader) io.Reader {
	return readerEOF{r}
}

type readerEOF struct {
	io.Reader
}

func (x readerEOF) Read(p []byte) (n int, err error) {
	n, err = x.Reader.Read(p)
	if n == 0 && err == nil {
		err = io.EOF
	}
	return
}

// ReadCloserEOF returns an io.ReadCloser which will return an io.EOF error as
// soon as the reader r is empty.
func ReadCloserEOF(r io.ReadCloser) io.ReadCloser {
	return readCloserEOF{r}
}

type readCloserEOF struct {
	io.ReadCloser
}

func (x readCloserEOF) Read(p []byte) (n int, err error) {
	n, err = x.ReadCloser.Read(p)
	if n == 0 && err == nil {
		err = io.EOF
	}
	return
}

// ReaderNopCloser
func ReaderNopCloser(r io.Reader) io.ReadCloser {
	return readerNopCloser{r}
}

type readerNopCloser struct {
	io.Reader
}

func (nr readerNopCloser) Close() error {
	return nil
}

// ReaderEOFNopCloser
func ReaderEOFNopCloser(r io.Reader) io.ReadCloser {
	return ReaderNopCloser(ReaderEOF(r))
}
