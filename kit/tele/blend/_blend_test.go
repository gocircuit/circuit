// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package blend

import (
	"encoding/gob"
	"os"
	"testing"
	"time"

	_ "github.com/hoijui/circuit/kit/debug/ctrlc"
	"github.com/hoijui/circuit/kit/tele/chain"
	"github.com/hoijui/circuit/kit/tele/codec"
	"github.com/hoijui/circuit/kit/tele/faithful"
	"github.com/hoijui/circuit/kit/tele/sandbox"
	"github.com/hoijui/circuit/kit/tele/trace"
)

type testMsg struct {
	Carry int
}

func init() {
	gob.Register(&testMsg{})
}

func failNow() {
	os.Exit(1)
}

const testN = 100

func TestCodec(t *testing.T) {

	f := trace.NewFrame()
	sx := sandbox.NewUnreliableTransport(f.Refine("sandbox"), 5, 0, time.Second/3, time.Second/3)
	hx := chain.NewTransport(f.Refine("chain"), sx)
	fx := faithful.NewTransport(f.Refine("faithful"), hx)
	cx := codec.NewTransport(fx, codec.GobCodec{})
	bx := NewTransport(f.Refine("session"), cx)

	// Sync
	ya, yb := make(chan int), make(chan int)

	// Accepter
	go func() {
		as := bx.Listen(sandbox.Addr("@")).AcceptSession()
		for i := 0; i < testN; i++ {
			go testAcceptConn(t, as, ya, yb)
		}
	}()

	// Dialer
	ds := bx.DialSession(sandbox.Addr("@"), nil)
	go func() {
		for i := 0; i < testN; i++ {
			go testDialConn(t, ds, ya, yb)
		}
	}()

	for i := 0; i < testN; i++ {
		<-yb
	}
}

func testDialConn(t *testing.T, ds *DialSession, ya, yb chan int) {
	<-ya
	conn := ds.Dial()
	if err := conn.Write(&testMsg{77}); err != nil {
		t.Errorf("write (%s)", err)
		failNow()
	}
	if err := conn.Close(); err != nil {
		t.Errorf("close (%s)", err)
		failNow()
	}
}

func testAcceptConn(t *testing.T, as *AcceptSession, ya, yb chan int) {
	ya <- 1
	conn := as.Accept()
	msg, err := conn.Read()
	if err != nil {
		t.Errorf("read (%s)", err)
		failNow()
	}
	if msg.(*testMsg).Carry != 77 {
		t.Errorf("check")
		failNow()
	}
	conn.Close()
	yb <- 1
}
