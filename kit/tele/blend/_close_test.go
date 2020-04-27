// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package blend

import (
	"net"
	"testing"
	"time"

	_ "github.com/hoijui/circuit/kit/debug/ctrlc"
	"github.com/hoijui/circuit/kit/tele/chain"
	"github.com/hoijui/circuit/kit/tele/codec"
	"github.com/hoijui/circuit/kit/tele/faithful"
	"github.com/hoijui/circuit/kit/tele/tcp"
	"github.com/hoijui/circuit/kit/tele/trace"
)

func TestClosure(t *testing.T) {
	const testK = 1

	f := trace.NewFrame()
	sx := tcp.Transport
	hx := chain.NewTransport(f.Refine("chain"), sx)
	fx := faithful.NewTransport(f.Refine("faithful"), hx)
	cx := codec.NewTransport(fx, codec.GobCodec{})
	bx := NewTransport(f.Refine("session"), cx)

	// Sync
	ya, yb := make(chan int), make(chan int)

	a, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:40090")

	// Accepter
	go func() {
		as := bx.Listen(a).AcceptSession()
		for i := 0; i < testK; i++ {
			go testAcceptConn(t, as, ya, yb)
		}
	}()

	// Dialer
	ds := bx.DialSession(a, nil)
	go func() {
		for i := 0; i < testK; i++ {
			go testDialConn(t, ds, ya, yb)
		}
	}()

	for i := 0; i < testK; i++ {
		<-yb
	}

	ds.Close()
	println("hold a minute...")
	time.Sleep(time.Minute)
	println("now check that the test process has no open tcp connections...")
	time.Sleep(time.Minute)
	println("great.")
}
