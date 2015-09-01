// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package assemble

import (
	"log"
	"net"

	"github.com/gocircuit/circuit/kit/xor"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
)

type Assembler struct {
	focus     xor.Key
	addr      n.Addr       // our circuit address
	multicast *net.UDPAddr // udp multicast address
}

func NewAssembler(addr n.Addr, multicast *net.UDPAddr) *Assembler {
	return &Assembler{
		focus:     xor.ChooseKey(),
		addr:      addr,
		multicast: multicast,
	}
}

func (a *Assembler) scatter(origin string) {
	msg := &TraceMsg{
		Origin: origin,
		Addr:   a.addr.String(),
	}
	scatter := NewScatter(a.multicast, a.focus, msg.Encode())
	scatter.Scatter() // send off a sequence of messages announcing our presnence over time
}

type JoinFunc func(n.Addr)

func (a *Assembler) AssembleServer(joinServer JoinFunc) {
	go a.scatter("server")
	go func() {
		gather := NewGather(a.multicast)
		for {
			_, payload := gather.Gather()
			trace, err := Decode(payload)
			if err != nil {
				log.Printf("Unrecognized trace message (%v)", err)
				continue
			}
			joinAddr, err := n.ParseAddr(trace.Addr)
			if err != nil {
				log.Printf("Trace origin address not parsing (%v)", err)
				continue
			}
			switch trace.Origin {
			case "server":
				joinServer(joinAddr)
			case "client":
				joinClient(a.addr, joinAddr)
			}
		}
	}()
}

func joinClient(serverAddr, clientAddr n.Addr) {
	x, err := circuit.TryDial(clientAddr, "dialback")
	if err != nil {
		return
	}
	y := YDialBack{x}
	y.OfferAddr(serverAddr)
}

func (a *Assembler) AssembleClient() n.Addr { // XXX: Clients should get more than one offering.
	d, xd := NewDialBack()
	circuit.Listen("dialback", xd)
	go a.scatter("client")
	return d.ObtainAddr()
}
