// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package interruptible

import (
	"bufio"
	"io"
	"runtime"
)

//
func Pipe() (Reader, Writer) {
	//
	ch := make(chan []byte)
	w := &writer{}
	w.w.ch = ch
	//
	abort := make(chan struct{})
	w.w.abort = abort
	w.a.abort = abort
	//
	r := &reader{w: w}
	r.r.ch = ch
	//
	runtime.SetFinalizer(r, func(r2 *reader) { r2.Close() })
	runtime.SetFinalizer(w, func(w2 *writer) { w2.Close() })
	//
	return r, w
}

//  ww<–pipe–>wr <–copy–> rw<–pipe–>rr
func BufferPipe(n int) (Reader, Writer) {
	x1, x0 := Pipe()
	x3, x2 := Pipe()
	go func() {
		bx2 := bufio.NewWriterSize(x2, n)
		io.Copy(bx2, x1)
		bx2.Flush()
		x2.Close()
	}()
	return x3, x0
}
