// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package bsg

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"runtime"

	"github.com/gocircuit/circuit/kit/xor"
)

type Gather struct {
	addr *net.UDPAddr // udp multicast address for discovery
	conn *net.UDPConn
}

func NewGather(addr *net.UDPAddr) *Gather {
	g := &Gather{
		addr: addr,
	}
	var err error
	if g.conn, err = net.ListenMulticastUDP("udp", nil, addr); err != nil {
		log.Printf("problem listening to udp mulsticast: %v", err)
		os.Exit(1)
	}
	runtime.SetFinalizer(g, 
		func(g2 *Gather) {
			g2.conn.Close()
		},
	)
	return g
}

func (s *Gather) Gather()  (xor.Key, []byte) {
	buf := make([]byte, 7e3)
	for {
		n, _, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		var msg Msg
		if err = json.Unmarshal(buf[:n], &msg); err != nil {
			continue // malformed invitation
		}
		return msg.Key, msg.Payload
	}
}

type GatherLens struct {
	gather *Gather
	lens *Lens
}

func NewGatherLens(addr *net.UDPAddr, focus xor.Key, k int) *GatherLens {
	return &GatherLens{
		gather: NewGather(addr),
		lens: NewLens(focus, k),
	}
}

func (s *GatherLens) Gather()  (xor.Key, []byte) {
	for {
		key, payload := s.gather.Gather()
		if s.lens.Remember(key) {
			return key, payload
		}
	}
}

func (s *GatherLens) Clear() {
	s.lens.Clear()
}
