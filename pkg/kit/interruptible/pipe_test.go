// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package interruptible

import (
	"io"
	"math/rand"
	"testing"
)

func TestPipe(t *testing.T) {
	r, w := Pipe()
	testPipe(r, w, t)
}

func TestBuffer(t *testing.T) {
	r, w := BufferPipe(1000)
	testPipe(r, w, t)
}

func testPipe(r io.ReadCloser, w io.WriteCloser, t *testing.T) {
	const N = 2200
	ch := make(chan int)
	data := make([]byte, N)
	for i := range data {
		data[i] = byte(rand.Int())
	}
	go func() { // write goroutine
		m := 0
		x := data
		for len(x) > 0 {
			n, err := w.Write(x)
			if err != nil {
				t.Fatalf("write: %v", err)
			}
			x = x[n:]
			m += n
		}
		w.Close()
		ch <- 1
	}()
	y := make([]byte, N)
	go func() { // read goroutine
		z := y
		m := 0
		for len(z) > 0 {
			n, err := r.Read(z)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			z = z[n:]
			m += n
		}
		ch <- 1
	}()
	<-ch
	<-ch
	for i := range data {
		if data[i] != y[i] {
			t.Fatalf("index %d differs: %v vs %v", i, data[i], y[i])
		}
	}
}
