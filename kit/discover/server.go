// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package discover

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	"github.com/gocircuit/circuit/kit/xor"
)

// circuit start -a :7711 -discover 228.8.8.8:8822

// Server is a network server for the beacon discovery protocol.
type Server struct {
	addr *net.UDPAddr // udp multicast address for discovery
	payload []byte // payload (circuit address) that we are advertising to the broadcast channel
	family *family
}

type InviteMsg struct {
	Payload []byte
}

// addr is a multicast address.
func New(addr *net.UDPAddr, payload []byte) (*Server, <-chan []byte) {
	ch := make(chan []byte)
	s := &Server{
		addr: addr,
		payload: payload,
		family: newFamily(xor.HashBytes(payload), 2),
	}
	s.Invite()
	go s.accept(ch)
	return s, ch
}

// Invite emits a burst of broadcasts, announcing this node, and resets the node's memory
func (s *Server) Invite() {
	s.family.Clear()
	go func() {
		conn, err := net.DialUDP("udp", nil, s.addr)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		msg := &InviteMsg{s.payload}
		buf, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}
		dur := time.Second
		for i := 0; i < 10; i++ {
			if _, err = conn.Write(buf); err != nil {
				log.Printf("invitation error: " + err.Error())
			}
			time.Sleep(dur)
			dur = (dur * 7) / 5
		}
	}()
}

// accept listens to broadcasts and chooses to join some of the newcomers, using an XOR-metric choice rule.
func (s *Server) accept(ch chan<- []byte) {
	conn, err := net.ListenMulticastUDP("udp", nil, s.addr)
	if err != nil {
		log.Printf("problem listening to udp mulsticast: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	buf := make([]byte, 7e3)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		var invite InviteMsg
		if err = json.Unmarshal(buf[:n], &invite); err != nil {
			continue // malformed invitation
		}
		if s.family.Remember(xor.HashBytes(invite.Payload)) {
			ch <- invite.Payload
		}
	}
}
