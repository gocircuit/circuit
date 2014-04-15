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
	runtime.SetFinalizer(w, func(w2 *writer) { w2.Close() })
	//
	return r, w
}

//  ww<–pipe–>wr <–copy–> rw<–pipe–>rr
func BufferPipe(n int) (Reader, Writer) {
	wr, ww := Pipe()
	rr, rw := Pipe()
	go func() {
		brw := bufio.NewWriterSize(rw, n)
		io.Copy(brw, wr)
		brw.Flush()
		rw.Close()
	}()
	return rr, ww
}
