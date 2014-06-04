// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package discover

import (
	"encoding/json"
	"net"
	"time"

	"github.com/gocircuit/circuit/kit/xor"
)

// circuit start -a :7711 -discover 4242

// Server is a network server for the beacon discovery protocol.
type Server struct {
	port int // udp broadcast port for discovery
	payload []byte // payload (circuit address) that we are advertising to the broadcast channel
	family *family
}

type InviteMsg struct {
	Payload []byte
}

func New(port int, payload []byte) (*Server, <-chan []byte) {
	ch := make(chan []byte)
	s := &Server{
		port: port,
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
		conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{ IP: net.IPv4bcast, Port: s.port })
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		msg := &InviteMsg{s.payload}
		buf, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}
		for i := 0; i < 3; i++ {
			if _, err = conn.Write(buf); err != nil {
				panic("invitation error: " + err.Error())
			}
			time.Sleep(time.Second)
		}
	}()
}

// accept listens to broadcasts and chooses to join some of the newcomers, using an XOR-metric choice rule.
func (s *Server) accept(ch chan<- []byte) {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{ IP: net.IPv4zero, Port: s.port })
	if err != nil {
		panic(err)
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
