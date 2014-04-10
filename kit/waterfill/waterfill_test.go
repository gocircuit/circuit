// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package waterfill

import (
	"fmt"
	"testing"
)

type testBin int

func (p *testBin) Add() {
	(*p)++
}

func (p *testBin) Less(fb FillBin) bool {
	return *p < *(fb.(*testBin))
}

func (p *testBin) String() string {
	return fmt.Sprintf("%02d", *p)
}

func TestFill(t *testing.T) {
	bin := make([]FillBin, 10)
	for i, _ := range bin {
		b := testBin(i * 2)
		bin[i] = &b
	}
	f := NewFill(bin)
	for i := 0; i < 30; i++ {
		println(f.String())
		f.Add()
	}
}
