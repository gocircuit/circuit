// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

// This is a circuit application that starts a node.js key/value service backed by a MySQL server.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"

	"github.com/gocircuit/circuit/client"
)

var flagAddr = flag.String("addr", "", "circuit server address, looks like circuit://...")

func fatalf(format string, arg ...interface{}) {
	println(fmt.Sprintf(format, arg...))
	os.Exit(1)
}

func connect(addr string) *client.Client {
	defer func() {
		if r := recover(); r != nil {
			fatalf("could not connect: %v", r)
		}
	}()
	// Connect the client to a circuit server
	return client.Dial(addr, nil)
}

func pickHosts(c *client.Client) (mysqlHost, nodejsHost client.Anchor) {
	defer func() {
		if recover() != nil {
			fatalf("client connection lost")
		}
	}()
	for _, s := range c.View() {
		return s, s // ???
	}
	fatalf("no available circuit hosts")
	return nil, nil
}

// runShell executes the shell command on the given host,
// waits until the command completes and returns its output
// as a string. The error value is non-nil if the process exited in error.
func runShell(host client.Anchor, shcmd string) (string, error) {
	defer func() {
		if recover() != nil {
			fatalf("connection to host lost")
		}
	}()
	job := host.Walk([]string{"shelljob", strconv.Itoa(rand.Int())})
	proc, _ := job.MakeProc(client.Cmd{
		Path:  "/bin/sh",
		Dir:   "/tmp",
		Args:  []string{"-c", shcmd},
		Scrub: true,
	})
	proc.Stdin().Close()  // Must close the standard input of the shell process.
	proc.Stderr().Close() // Close to indicate discarding standard error
	var buf bytes.Buffer
	io.Copy(&buf, proc.Stdout())
	stat, _ := proc.Wait()
	return buf.String(), stat.Exit
}

func main() {
	flag.Parse()

	c := connect(*flagAddr)

	mysqlHost, _ /*, nodejsHost*/ := pickHosts(c)

	out, err := runShell(mysqlHost, "ls -l /")
	if err != nil {
		println("err:", err)
	} else {
		println("ok:", out)
	}
}
