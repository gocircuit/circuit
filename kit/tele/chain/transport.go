// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chain

import (
	"net"
	"time"

	"github.com/gocircuit/circuit/kit/tele/trace"
)

type Transport struct {
	frame   trace.Frame
	stat    Stat
	carrier Carrier
}

func NewTransport(frame trace.Frame, carrier Carrier) *Transport {
	t := &Transport{
		frame:   frame,
		carrier: carrier,
	}
	t.frame.Bind(t)
	t.stat.Init()
	/*
		go func() {
			s := t.Subscribe()
			for {
				u := s.Wait()
				t.frame.Printf("Status: NDial=%d NAccept=%d", u.NDial, u.NAccept)
				//println(t.stat.Dump())
			}
		}()
	*/
	return t
}

func (t *Transport) Subscribe() *Subscriber {
	return t.stat.Subscribe()
}

const CarrierRedialTimeout = time.Second

// Dial creates a new chain connection to addr.
// Dial returns instanteneously (it does not wait on I/O operations) and always succeeds, returning a non-nil connection object.
func (t *Transport) Dial(addr net.Addr) *Conn {
	id := chooseChainID()
	t.stat.addDC(id, addr)
	dc := newDialConn(t.frame.Refine(addr.String()), id, addr,
		func() (net.Conn, error) {
			return t.carrier.Dial(addr)
		},
		func() {
			t.stat.scrubDC(id)
		},
	)
	return &dc.Conn
}

func (t *Transport) Listen(addr net.Addr) *Listener {
	return NewListener(t.frame.Refine("listener"), &t.stat, t.carrier, addr)
}
