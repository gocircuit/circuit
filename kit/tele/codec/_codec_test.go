// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package codec

import (
	"os"
	"testing"
	"time"

	_ "github.com/hoijui/circuit/kit/debug/ctrlc"
	"github.com/hoijui/circuit/kit/tele/chain"
	"github.com/hoijui/circuit/kit/tele/faithful"
	"github.com/hoijui/circuit/kit/tele/sandbox"
	"github.com/hoijui/circuit/kit/tele/trace"
)

type testMsg struct {
	Carry int
}

func failNow() {
	os.Exit(1)
}

const testN = 5

func TestCodec(t *testing.T) {

	// Transport
	f := trace.NewFrame()
	// Carrier
	sx := sandbox.NewRandomUnreliableTransport(f.Refine("sandbox"), 3, 3, time.Second/4, time.Second/4)
	// Chain
	hx := chain.NewTransport(f.Refine("chain"), sx)
	// Faithful
	fx := faithful.NewTransport(f.Refine("faithful"), hx)
	// Codec
	cx := NewTransport(fx, GobCodec{})

	// Sync
	y := make(chan int)

	// Accepter
	go func() {
		l := cx.Listen(sandbox.Addr("@"))
		for i := 0; i < testN; i++ {
			y <- 1
			conn := l.Accept()
			msg := &testMsg{}
			if err := conn.Read(msg); err != nil {
				t.Fatalf("read (%s)", err)
				failNow()
			}
			if msg.Carry != i {
				t.Fatalf("check")
				failNow()
			}
			f.Printf("READ %d/%d CLOSING", i+1, testN)
			conn.Close()
			f.Printf("READ %d/%d √", i+1, testN)
		}
		y <- 1
	}()

	// Dialer
	for i := 0; i < testN; i++ {
		<-y
		conn := cx.Dial(sandbox.Addr("@"))
		if err := conn.Write(&testMsg{i}); err != nil {
			t.Fatalf("write (%s)", err)
			failNow()
		}
		f.Printf("WRITE %d/%d CLOSING", i+1, testN)
		if err := conn.Close(); err != nil {
			t.Fatalf("close (%s)", err)
			failNow()
		}
		f.Printf("WRITE %d/%d √", i+1, testN)
	}
	<-y
}
