// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package sandbox

import (
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/gocircuit/circuit/kit/tele/trace"
)

const N = 50

func TestReliable(t *testing.T) {

	sandbox := NewReliableTransport(trace.NewFrame("TestReliable"))
	ready := make(chan int)

	go func() {
		l, err := sandbox.Listen(Addr("x"))
		if err != nil {
			t.Fatalf("listen (%s)", err)
		}
		for i := 0; i < N; i++ {
			ready <- 1
			//println("receiver", i)
			conn, err := l.Accept()
			if err != nil {
				t.Fatalf("accept (%s)", err)
			}
			r := make([]byte, 3)
			if _, err := conn.Read(r); err != nil {
				t.Fatalf("read (%s)", err)
			}
			if !reflect.DeepEqual([]byte{0, 1, 2}, r) {
				t.Fatalf("unexpected read")
			}
			var w []byte = []byte{0, 1, 2}
			if _, err := conn.Write(w); err != nil {
				t.Fatalf("write (%s)", err)
			}
			if err := conn.Close(); err != nil {
				t.Fatalf("close (%s)", err)
			}
		}
		ready <- 1
	}()

	for i := 0; i < N; i++ {
		<-ready
		//println("sender", i)
		conn, err := sandbox.Dial(Addr("x"))
		if err != nil {
			t.Fatalf("dial (%s)", err)
		}
		var w []byte = []byte{0, 1, 2}
		if _, err := conn.Write(w); err != nil {
			t.Fatalf("write (%s)", err)
		}
		r := make([]byte, 3)
		if _, err := conn.Read(r); err != nil {
			t.Fatalf("read (%s)", err)
		}
		if err := conn.Close(); err != nil {
			t.Fatalf("close (%s)", err)
		}
	}
	<-ready
}

func TestHalfConn(t *testing.T) {
	ready := make(chan int)
	f := trace.NewFrame("TestHalfConn")
	p, q := NewPipe(f.Refine("@r"), f.Refine("@w"), Addr("@r"), Addr("@w"))
	// Read goroutine
	go func() {
		defer func() {
			close(ready)
		}()
		b, err := ioutil.ReadAll(q)
		if err != nil {
			log.Printf("read all (%s)", err)
			t.Fatalf("read all (%s)", err)
		}
		if len(b) != N {
			t.Fatalf("bad length")
		}
		for i := 0; i < N; i++ {
			if b[i] != byte(i) {
				t.Fatalf("unexpected byte @%d", i)
			}
		}
	}()
	// Write logic
	for i := 0; i < N; i++ {
		if n, err := p.Write([]byte{byte(i)}); err != nil || n != 1 {
			t.Fatalf("write %d (%v)", n, err)
		}
	}
	if err := p.Close(); err != nil {
		t.Fatalf("close (%s)", err)
	}
	<-ready
}

func TestUnreliable(t *testing.T) {
	f := trace.NewFrame("TestDropTail")
	q, p := NewSievePipe(f.Refine("@r"), f.Refine("@w"), Addr("@r"), Addr("@w"), N, N, 0, 0)
	// Read goroutine
	ready := make(chan int)
	go func() {
		defer func() {
			close(ready)
		}()
		b, err := ioutil.ReadAll(q)
		// We get io.ErrUnexpectedEOF because the closing eof{} is not transmitted properly
		if err != io.ErrUnexpectedEOF {
			t.Fatalf("read all (%s)", err)
		}
		if len(b) != N {
			t.Fatalf("bad length")
		}
		for i := 0; i < N; i++ {
			if b[i] != byte(i) {
				t.Fatalf("unexpected byte @%d", i)
			}
		}
	}()
	// The following writes go through
	for i := 0; i < N; i++ {
		if n, err := p.Write([]byte{byte(i)}); err != nil || n != 1 {
			t.Fatalf("write %d (%v)", n, err)
		}
	}
	// The following writes are all dropped silently
	for i := 0; i+1 < N; i++ {
		if n, err := p.Write([]byte{byte(i)}); err != nil || n != 1 {
			t.Fatalf("write %d (%v)", n, err)
		}
	}
	// The last write is dropped quietly and also returns the unexpected EOF
	/*
		if n, err := p.Write([]byte{N - 1}); err != io.ErrUnexpectedEOF || n != 1 {
			t.Fatalf("write expects error (%v), got (%v); expects len %d, got %d", io.ErrUnexpectedEOF, err, 1, n)
		}
	*/
	if err := p.Close(); err != nil {
		t.Fatalf("close (%s)", err)
	}
	<-ready
}
