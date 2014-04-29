// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestChanSel(t *testing.T) {
	slash := os.Getenv("LOOPBACK")
	//
	println("make chan")
	charlie := path.Join(slash, "chan", "charlie")
	os.MkdirAll(charlie, 0777)
	ioutil.WriteFile(path.Join(charlie, "cap"), []byte("0"), 0)
	//
	println("make select")
	sarah := path.Join(slash, "select", "sarah")
	os.MkdirAll(sarah, 0777)
	ioutil.WriteFile(
		path.Join(sarah, "select"), 
		[]byte(fmt.Sprintf(`[{"op": "r", "file": "%s"}]`, path.Join(charlie, "recv"))),
		0,
	)
	//
	//println("send on chan")
	//ioutil.WriteFile(path.Join(charlie, "send"), []byte("¡hello world!"), 0)
	//
	// BUG: cannot ls in charlie directory at this point
	//
	println("ok check")
	//
	<-(chan int)(nil)
}