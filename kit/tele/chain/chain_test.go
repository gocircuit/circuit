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
	"time"

	_ "github.com/gocircuit/circuit/kit/debug/ctrlc"
	"github.com/gocircuit/circuit/kit/tele/sandbox"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

type testKind int

const (
	testDRAW = testKind(iota) // Dialer Reads, Accepter Writes
	testDWAR                  // Dialer Writers, Accepter Reads
)

type connMode struct {
	Kind       testKind
	Random     bool
	NOK, NDrop int
	Expect     []byte
}

func (mode connMode) NWrite() int {
	return int(mode.Expect[len(mode.Expect)-1] + 1)
}

var testModes = []connMode{
	// Deterministic tests
	// DWAR
	{testDWAR, false, 2, 4, []byte{0, 5, 10}},                  // dial, 0, (1), (2), (3), (4), dial, 5, (6), (7), (8), (9), dial, 10
	{testDWAR, false, 4, 0, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8}}, // dial, 0, 1, 2, dial, 3, 4, 5, dial, 6, 7, 8
	{testDWAR, false, 3, 2, []byte{0, 1, 4, 5, 8, 9, 12, 13}},  // dial, 0, 1, (2), (3), dial, 4, 5, (6), (7), dial, 8, 9, (10), (11), dial, 12, 13
	// DRAW
	{testDRAW, false, 4, 0, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}, // dial, 0, 1, 2, (STITCH), dial, 3, 4, 5, (STITCH), ...
	{testDRAW, false, 2, 0, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	{testDRAW, false, 3, 2, []byte{0, 1, 2, 5, 6, 7, 10, 11, 12}},
	{testDRAW, false, 1, 2, []byte{0, 3, 6, 9}},
	// Random tests
	{testDWAR, true, 5, 5, []byte{100}}, // Dialer writes 0, 1, 2, ..., 99 and carrier connections are broken randomly
	{testDRAW, true, 5, 5, []byte{100}}, // Accepter writes 0, 1, 2, ..., 99 and carrier connections are broken randomly
}

func TestConn(t *testing.T) {
	for _, m := range testModes {
		testConn(t, m)
	}
}

func failNow() {
	// Wait a bit to prevent race between log.Printfs and the panic's println.
	time.Sleep(time.Second / 4)
	panic("force fail")
}

func testConn(t *testing.T, mode connMode) {
	// Transport
	frame := trace.NewFrame()
	var x Carrier
	if mode.Random {
		x = sandbox.NewRandomUnreliableTransport(frame.Refine("sandbox"), mode.NOK, mode.NDrop, 0, 0)
	} else {
		x = sandbox.NewUnreliableTransport(frame.Refine("sandbox"), mode.NOK, mode.NDrop, 0, 0)
	}
	// Signalling
	ready := make(chan int, 2)
	feedfwd := make(chan byte, mode.NWrite()+1)

	// Accepter endpoint
	go func() {
		defer func() {
			ready <- 1 // SYNC: Notify that accepter-side logic is done
		}()

		l := NewListener(frame.Refine("chain", "listener"), x, sandbox.Addr(""))
		ready <- 1 // SYNC: Notify that listener is accepting
		switch mode.Kind {
		case testDWAR:
			testRead(t, l.Accept(), feedfwd)
		case testDRAW:
			testWrite(t, l.Accept(), mode, feedfwd)
		default:
			panic("u")
		}
	}()

	// Dialer endpoint
	<-ready // SYNC: Wait for listener to start accepting
	d := NewDialer(frame.Refine("chain", "dialer"), x)
	conn := d.Dial(sandbox.Addr(""))
	switch mode.Kind {
	case testDWAR:
		testWrite(t, conn, mode, feedfwd)
	case testDRAW:
		testRead(t, conn, feedfwd)
	default:
		panic("u")
	}
	<-ready // SYNC: Wait for accepter goroutine to complete
}

func testRead(t *testing.T, conn *Conn, feedfwd <-chan byte) {
	var i int
	ch := make(chan []byte, 100)
	// Read loop
	go func() {
		var v []byte
		var err error
		for {
			v, err = conn.Read()
			if err != nil {
				if v != nil {
					panic("eh")
				}
				// If not a stitching error, connection killed
				if IsStitch(err) == nil {
					return
				}
				continue
			}
			ch <- v
		}
	}()
	// Feed forward loop
	for u := range feedfwd {
		v := <-ch
		log.Printf("READ #%d = 0x%x", i+1, v[0])
		if !reflect.DeepEqual(v, []byte{byte(u)}) {
			t.Errorf("expecting %#v, got %#v", []byte{byte(u)}, v)
			log.Printf("expecting %#v, got %#v", []byte{byte(u)}, v)
			failNow()
		}
		log.Printf("CHKD #%d = 0x%x", i+1, v[0])
		i++
	}
	conn.Kill()
	log.Println("READ KILLED")
}

func testWrite(t *testing.T, conn *Conn, mode connMode, feedfwd chan<- byte) {
	n := mode.NWrite()
	ready := make(chan struct{})

	// If not random, send the expected reception data to remote side
	if !mode.Random {
		go func() {
			for _, b := range mode.Expect {
				feedfwd <- b
			}
			close(ready)
		}()
	}
	// Write
	var i int
	for i < n {
		_, err := conn.Read()
		w := IsStitch(err)
		if w == nil {
			log.Printf("not stitch read on writer side")
			failNow()
		}
		for i < n {
			err := w.Write([]byte{byte(i)})
			i++
			// If a random test, tell the receiver out-of-band what they are supposed to receive
			if w.Delivered() {
				if mode.Random {
					feedfwd <- byte(i - 1)
				}
				log.Printf("WROTE %d/%d = 0x%x", i, n, i-1)
			}
			if err != nil {
				// Get a new writer
				break
			}
		}
	}

	if mode.Random {
		close(ready)
	}
	<-ready // SYNC: Wait until all feedback has been sent to the reader
	conn.Kill()
	log.Println("WRITE KILLED")
	close(feedfwd) // SYNC
}
