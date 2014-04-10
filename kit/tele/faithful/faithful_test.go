// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package faithful

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/gocircuit/circuit/kit/tele/chain"
	"github.com/gocircuit/circuit/kit/tele/sandbox"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

const N = 20

type testKind int

const (
	testDRAW = testKind(iota) // Dialer Reads, Accepter Writes
	testDWAR                  // Dialer Writers, Accepter Reads
)

type testMode struct {
	Kind       testKind
	Random     bool
	NOK, NDrop int
}

var testModes = []testMode{
	// DWAR
	{testDWAR, false, 3, 0},
	{testDWAR, false, 3, 1},
	{testDWAR, false, 4, 0},
	{testDWAR, false, 7, 2},
	{testDWAR, false, 3, 3},
	{testDWAR, false, 3, 7},
	// DRAW
	{testDRAW, false, 3, 0},
	{testDRAW, false, 3, 1},
	{testDRAW, false, 4, 0},
	{testDRAW, false, 7, 2},
	{testDRAW, false, 3, 3},
	// Random tests
	{testDWAR, true, 5, 5},
	{testDRAW, true, 5, 5},
}

func TestConn(t *testing.T) {
	for _, mode := range testModes {
		log.Printf("------------------- STARING TEST %#v ----------------------", mode)
		testConn(t, mode)
	}
}

const exp = time.Second * 2 // After expiration of inactivty on the send side, the dropTail will break the connection

type aborter interface {
	Abort()
}

var (
	abrt__ sync.Mutex
	abrt   func()
)

func setabrt(a func()) {
	abrt__.Lock()
	defer abrt__.Unlock()
	abrt = a
}

func init() {
	go func() {
		ch := make(chan os.Signal, 2)
		signal.Notify(ch, os.Interrupt)
		var i int
		for _ = range ch {
			i++
			if i == 1 {
				abrt__.Lock()
				a := abrt
				abrt__.Unlock()
				a()
				continue
			}
			panic("ctrlc")
		}
	}()
}

func testConn(t *testing.T, mode testMode) {
	f := trace.NewFrame("testConn")
	f.Bind(&f)
	fsx := f.Refine("sandbox")
	fsx.Bind(&fsx)

	var x *sandbox.Transport
	if mode.Random {
		x = sandbox.NewRandomUnreliableTransport(fsx, mode.NOK, mode.NDrop, exp, exp)
	} else {
		x = sandbox.NewUnreliableTransport(fsx, mode.NOK, mode.NDrop, exp, exp)
	}

	ready := make(chan int, 3)

	// Accept side
	go func() {
		ax := NewTransport(f.Refine("faithful:a"), chain.NewTransport(f.Refine("chain:a"), x))
		l := ax.Listen(sandbox.Addr("@"))
		ready <- 1
		c := l.Accept()
		switch mode.Kind {
		case testDWAR:
			testRead(f, t, c, ready)
		case testDRAW:
			testWrite(f, t, c, ready)
		}
	}()

	// Dial side
	dx := NewTransport(f.Refine("faithful:d"), chain.NewTransport(f.Refine("chain:d"), x))
	<-ready
	c := dx.Dial(sandbox.Addr("@"))
	setabrt(func() {
		dbg := c.Debug().(*sandbox.DebugInfo)
		dbg.Out.Abort()
	})
	switch mode.Kind {
	case testDWAR:
		testWrite(f, t, c, ready)
	case testDRAW:
		testRead(f, t, c, ready)
	}
	_, _ = <-ready, <-ready // One from testRead and one from testWrite
}

func testRead(fr trace.Frame, t *testing.T, c *Conn, ready chan<- int) {
	defer func() {
		ready <- 1
	}()
	for i := 0; i < N; i++ {
		q, err := c.Read()
		if err != nil {
			t.Fatalf("read (%s)", err)
			failNow()
		}
		z := []byte{byte(i), byte(i + 1), byte(i + 2)}
		if !reflect.DeepEqual(q, z) {
			t.Fatalf("expecting %#v, got %#v", z, q)
			failNow()
		}
		fr.Printf("READ %d/%d", i+1, N)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("read-side close (%s)", err)
		failNow()
	}
	fr.Printf("CLOSED READ")
}

func testWrite(fr trace.Frame, t *testing.T, c *Conn, ready chan<- int) {
	defer func() {
		ready <- 1
	}()
	for i := 0; i < N; i++ {
		if err := c.Write([]byte{byte(i), byte(i + 1), byte(i + 2)}); err != nil {
			t.Errorf("write (%s)", err)
			failNow()
		}
		fr.Printf("WROTE %d/%d", i+1, N)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("write-side close (%s)", err)
		failNow()
	}
	fr.Printf("CLOSED WRITE")
}

func failNow() {
	os.Exit(1)
}
