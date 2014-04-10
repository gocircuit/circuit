// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package iomisc

import (
	"io"
	"os"
	"testing"
)

func TestPrefixReader(t *testing.T) {
	pr, w, err := os.Pipe()
	if err != nil {
		panic(err.Error())
	}
	r := PrefixReader("R: ", pr)
	go func() {
		for _, l := range lines {
			_, err := w.Write([]byte(l))
			if err != nil {
				t.Fatalf("write (%s)", err)
			}
		}
		w.Close()
	}()
	io.Copy(os.Stderr, r)
}

func TestPrefixWriter(t *testing.T) {
	w := PrefixWriter("W: ", os.Stderr)
	for _, l := range lines {
		_, err := w.Write([]byte(l))
		if err != nil {
			t.Fatalf("write (%s)", err)
		}
	}
}
