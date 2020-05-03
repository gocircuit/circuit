// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package shell

/*
import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestShell(t *testing.T) {
	sh := NewShellServer(".")
	stdin, stdout, stderr, err := sh.Shell()
	if err != nil {
		t.Fatalf("shell (%s)", err)
		panic("q")
	}
	y := make(chan int)
	println("•shell")
	go func() {
		n, err := io.Copy(os.Stdout, stdout)
		fmt.Printf("stdout<- %d %#v\n", n, err)
		y <- 1
	}()
	go func() {
		n, err := io.Copy(os.Stderr, stderr)
		fmt.Fprintf(os.Stderr, "stdout<- %d %#v\n", n, err)
		y <- 1
	}()
	go io.Copy(stdin, os.Stdin)
	println("•0th")
	<-y
	println("•1st")
	<-y
	println("•2nd")
}
*/
