// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package interruptible

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

//
type Reader interface {
	io.ReadCloser
	ReadIntr([]byte, Intr) (int, error)
}

//
type reader struct {
	w *writer
	r struct {
		ch <-chan []byte
		Mutex
		buf bytes.Buffer
	}
	s struct {
		sync.Mutex
		n      int64
		closed bool
	}
}

func (r *reader) Read(p []byte) (n int, err error) {
	return r.ReadIntr(p, nil)
}

func (r *reader) ReadIntr(p []byte, intr Intr) (n int, err error) {
	u := r.r.Lock(intr)
	if u == nil {
		return 0, ErrIntr
	}
	defer u.Unlock()
	//
	defer func() {
		r.s.Lock()
		defer r.s.Unlock()
		//
		r.s.n += int64(n)
		if err != nil {
			r.s.closed = true
		}
	}()
	//
	if r.r.buf.Len() > 0 {
		return r.r.buf.Read(p)
	}
	//
	r.s.Lock()
	closed := r.s.closed
	r.s.Unlock()
	if closed {
		return 0, io.ErrClosedPipe
	}
	//
	select {
	case block, ok := <-r.r.ch:
		if !ok {
			return 0, io.EOF
		}
		r.r.buf.Write(block)
		return r.r.buf.Read(p)
	case <-intr:
		// Next message is not extracted
		return 0, errors.New("no progress") // io.ErrNoProgress
	}
	panic(0)
}

func (r *reader) Close() error {
	r.w.Close()
	return nil
}

func (r *reader) Stat() (nrecv int64, closed bool) {
	r.s.Lock()
	defer r.s.Unlock()
	//
	nrecv, closed = r.s.n, r.s.closed
	return
}
