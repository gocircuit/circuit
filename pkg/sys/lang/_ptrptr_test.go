// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestPtrPtr(t *testing.T) {
	l1 := NewSandbox()
	r1 := New(l1, &testBoot{"π1"})

	l2 := NewSandbox()
	r2 := New(l2, &testBoot{"π2"})

	p2, err := r1.TryDial(l2.Addr())
	if err != nil {
		t.Fatalf("dial 1->2 (%s)", err)
	}

	p1, err := r2.TryDial(l1.Addr())
	if err != nil {
		t.Fatalf("dial 2->1 (%s)", err)
	}

	if p1.Call("Name")[0].(string) != "π1" {
		t.Errorf("return val 1")
	}

	if p2.Call("Name")[0].(string) != "π2" {
		t.Errorf("return val 2")
	}
	p2.Call("ReturnNilMap")
}
