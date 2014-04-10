// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"github.com/gocircuit/circuit/use/n"
	"testing"
)

const tN = 100

func TestSystem(t *testing.T) {
	fs := NewSystem()
	x1, x2 := n.ChooseWorkerID(), n.ChooseWorkerID()
	x3, x4 := n.ChooseWorkerID(), n.ChooseWorkerID()
	ch := make(chan int)
	go func() {
		for i := 0; i < tN; i++ {
			fs.Create([]string{"/a", "/a"}, x1)
			fs.Create([]string{"/a", "/a/b"}, x2)
			fs.Remove(x2)
			fs.Remove(x1)
		}
		ch <- 1
	}()
	go func() {
		for i := 0; i < tN; i++ {
			fs.Create([]string{"/a/b", "/x"}, x3)
			fs.Create([]string{"/a", "/x/y"}, x4)
			fs.Remove(x3)
			fs.Remove(x4)
		}
		ch <- 1
	}()
	<-ch
	<-ch
	println(fs.Dump())
}
