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
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/gocircuit/circuit/client"
)

var flagAddr = flag.String("addr", "", "circuit server address, looks like circuit://...")

func fatalf(format string, arg ...interface{}) {
	println(fmt.Sprintf(format, arg...))
	os.Exit(1)
}

// connect establishes a client connection to the circuit cluster (via the given circuit server address)
// and returns a connected client object.
func connect(addr string) *client.Client {
	defer func() {
		if r := recover(); r != nil {
			fatalf("could not connect: %v", r)
		}
	}()
	return client.Dial(addr, nil)
}

func pickHosts(c *client.Client, n int) (hosts []client.Anchor) {
	defer func() {
		if recover() != nil {
			fatalf("client connection lost")
		}
	}()
	for _, a := range c.View() {
		if len(hosts) >= n {
			break
		}
		hosts = append(hosts, a)
	}
	if len(hosts) != n {
		fatalf("not enough available hosts found")
	}
	return
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

func getDarwinHostIP(host client.Anchor) string {
	out, err := runShell(host, `ifconfig en0 | awk '/inet / {print $2}'`)
	if err != nil {
		fatalf("get ip error: %v", err)
	}
	out = strings.TrimSpace(out)
	if _, err := net.ResolveIPAddr("ip", out); err != nil {
		fatalf("ip %q unrecognizable: %v", out, err)
	}
	return out
}

func getUbuntuHostIP(host client.Anchor) string {
	out, err := runShell(host, `ifconfig eth0 | awk '/inet addr/ {split($2, a, ":"); print a[2] }'`)
	if err != nil {
		fatalf("get ip error: %v", err)
	}
	out = strings.TrimSpace(out)
	if _, err := net.ResolveIPAddr("ip", out); err != nil {
		fatalf("ip %q unrecognizable: %v", out, err)
	}
	return out
}

func startMysql(host client.Anchor) (ip, port string) {
	// ??
}

func main() {
	flag.Parse()

	c := connect(*flagAddr)

	host := pickHosts(c, 1) // ??

	mysqlIP, mysqlPort := startMysql(host[0])

	// nodejsIP, nodejsPort := startNodejs(host[1], mysqlIP, mysqlPort)

	// println(getDarwinHostIP(hosts[0]))
}
