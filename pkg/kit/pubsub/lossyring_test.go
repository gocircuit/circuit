// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package pubsub

import (
	"runtime"
	"testing"
)

func TestLossyRing(t *testing.T) {
	r := MakeLossyRing(3)
	if !r.Send(1) {
		t.Fatalf("x")
	}
	if !r.Send(2) {
		t.Fatalf("x")
	}
	if r.Len() != 2 {
		t.Fatal("len", r.Len())
	}
	if !r.Send(3) {
		t.Fatalf("x")
	}
	if r.Send(4) {
		t.Fatalf("x")
	}
	if v, ok := r.Recv(); !ok {
		t.Fatalf("x")
	} else {
		if w, ok := v.(Loss); !ok || w.Count != 2 {
			t.Fatalf("x")
		}
	}
	if v, ok := r.Recv(); !ok {
		t.Fatalf("x")
	} else {
		if w, ok := v.(int); !ok || w != 3 {
			t.Fatalf("x")
		}
	}
}

func smrz() []interface{} {
	return []interface{}{1,2,3}
}

func TestPubSub(t *testing.T) {
	ps := New("nm", smrz)
	ps.Source()
	ch := make(chan int)
	go func() {
		ps.Publish(4)
		ps.Close()
		ch <- 1
	}()
	go func() {
		s := ps.Subscribe()
		s.Peek()
		// for _, ok := s.Consume(); ok; _, ok = s.Consume() {}
		s.Consume()
		s.Scrub()
		ch <- 1
	}()
	<-ch
	<-ch
	runtime.GC()
}
