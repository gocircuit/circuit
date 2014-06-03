// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tele

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/gocircuit/circuit/use/n"
)

type testMsg struct{}

func init() {
	gob.Register(&testMsg{})
}

func TestTele(t *testing.T) {
	sys := &System{}
	y, z := make(chan int), make(chan int)
	// Listener
	go func() {
		x := sys.NewTransport("listener")
		l := x.Listen(MustParseNetAddr("0.0.0.0:44111"))
		z <- 1
		conn := l.Accept()
		msg, err := conn.Read()
		if err != nil {
			failnow("listener read (%s)", err)
		}
		if _, ok := msg.(*testMsg); !ok {
			failnow("listener message type")
		}
		if err = conn.Close(); err != nil {
			failnow("listener close (%s)", err)
		}
		y <- 1
	}()
	// Dialer
	go func() {
		x := sys.NewTransport("dialer")
		<-z
		conn, err := x.Dial(&Addr{
			ID:  n.WorkerID("listener"),
			PID: os.Getpid(),
			TCP: MustParseNetAddr("127.0.0.1:44111").(*net.TCPAddr),
		})
		if err != nil {
			failnow("dialer dial (%s)", err)
		}
		if err = conn.Write(&testMsg{}); err != nil {
			failnow("dialer write (%s)", err)
		}
		if err = conn.Close(); err != nil {
			failnow("dialer close (%s)", err)
		}
		y <- 1
	}()
	<-y
	<-y
}

func failnow(format string, v ...interface{}) {
	println(fmt.Sprintf(format, v...))
	os.Exit(1)
}
