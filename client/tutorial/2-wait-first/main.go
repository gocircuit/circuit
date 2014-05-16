// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"os"
	"strconv"

	"github.com/gocircuit/circuit/client"
)

const n = 5

func pick(c *client.Client) client.Anchor {
	for _, r := range c.View() {
		return r
	}
	panic(0)
}

func main() {
	c := client.Dial(os.Args[1]) // argument is the url of a circuit server
	ch := make(chan int)
	for i := 0; i < n; i++ {
		cmd := client.Cmd{
			Path: "/bin/sleep",
			Args: []string{strconv.Itoa(3+i*3)},
		}
		i_ := i
		go func() {
			t := pick(c).Walk([]string{"wait_first", strconv.Itoa(i_)})
			p, _ := t.MakeProc(cmd)
			p.Stdin().Close()
			p.Wait()
			ch <- 1
			t.Scrub()
			println("process", i_+1, "done")
		}()
	}
	<-ch
	println("one done.")
}
