// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package interruptible

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
	return r, w
}
