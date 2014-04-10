// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chain

import (
	"log"
	"reflect"
	"testing"

	_ "github.com/gocircuit/circuit/kit/debug/ctrlc"
	"github.com/gocircuit/circuit/kit/tele/tcp"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

func TestConnOverTCP(t *testing.T) {
	frame := trace.NewFrame()
	x := tcp.Transport

	ready := make(chan int, 2)
	sent, recv := make(map[byte]struct{}), make(map[byte]struct{})

	// Accepter endpoint
	go func() {
		l := NewListener(frame.Refine("listener"), x, tcp.Addr(":17222"))
		ready <- 1 // SYNC: Notify that listener is accepting
		testGreedyRead(t, l.Accept(), recv)
		ready <- 1
	}()

	// Dialer endpoint
	<-ready // SYNC: Wait for listener to start accepting
	d := NewDialer(frame.Refine("dialer"), x)
	conn := d.Dial(tcp.Addr("localhost:17222"))
	testGreedyWrite(t, conn, sent)
	<-ready // SYNC: Wait for accepter goroutine to complete

	// Make sure all marked writes have been received
	if !reflect.DeepEqual(sent, recv) {
		t.Errorf("expected %#v, got %#v", sent, recv)
		failNow()
	}
}

func testGreedyRead(t *testing.T, conn *Conn, recv map[byte]struct{}) {
	var i int
	for i < testN {
		v, err := conn.Read()
		if err != nil {
			if v != nil {
				panic("eh")
			}
			// If not a stitching error, connection killed
			if IsStitch(err) == nil {
				return
			}
			log.Println("READ STITCH")
			continue
		}
		log.Printf("READ %d", v[0])
		recv[byte(v[0])] = struct{}{}
		i++
	}
	conn.Kill()
	log.Println("READ KILLED")
}

const testN = 5

func testGreedyWrite(t *testing.T, conn *Conn, sent map[byte]struct{}) {
	var i int
	for i < testN {
		_, err := conn.Read()
		w := IsStitch(err)
		if w == nil {
			log.Printf("non-stitch read on writer side")
			failNow()
		}
		for i < testN {
			log.Printf("WRITE %d", i)
			err := w.Write([]byte{byte(i)})
			if err != nil {
				t.Errorf("write (%s)", err)
				failNow()
			}
			sent[byte(i)] = struct{}{}
			i++
		}
	}
	conn.Kill()
	log.Println("WRITE KILLED")
}
