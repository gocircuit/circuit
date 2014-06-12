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
	"time"

	"github.com/gocircuit/circuit/kit/xor"
)

type Scatter struct {
	addr *net.UDPAddr // udp multicast address
	key xor.Key
	payload []byte
}

type Msg struct {
	Key xor.Key
	Payload []byte
}

// addr is a udp multicast address.
func NewScatter(addr *net.UDPAddr, key xor.Key, payload []byte) *Scatter {
	return &Scatter{
		addr: addr,
		key: key,
		payload: payload,
	}
}

func (s *Scatter) Scatter() {
	conn, err := net.DialUDP("udp", nil, s.addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	msg := &Msg{
		Key: s.key,
		Payload: s.payload,
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	dur := time.Second
	for i := 0; i < 10; i++ {
		if _, err = conn.Write(buf); err != nil {
			log.Printf("multicast scatter error: " + err.Error())
		}
		time.Sleep(dur)
		dur = (dur * 7) / 5
	}
}
