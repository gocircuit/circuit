// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package dns

import (
	"net"
	"sync"

	"github.com/gocircuit/circuit/github.com/miekg/dns"
	"github.com/gocircuit/circuit/use/circuit"
)

type Nameserver interface {
	Scrub()
	Set(rr string) error
	Unset(name string)
	Peek() Stat
	X() circuit.X
}

type nameserver struct {
	sync.Mutex
	server *dns.Server
	addr net.Addr
	rr map[string][]dns.RR // name -> rr
}

func MakeNameserver() (_ Nameserver, err error) {
	ns := &nameserver{
		rr: make(map[string][]dns.RR),
	}
	if err = ns.startUdpServer(); err != nil {
		return nil, err
	}
	return ns, nil
}

func (ns *nameserver) startUdpServer() error {
	pc, err := net.ListenPacket("udp", "") // empty-string address picks an available port on 0.0.0.0
	if err != nil {
		return err
	}
	ns.server = &dns.Server{
		PacketConn: pc,
		Handler: ns,
	}
	ns.addr = pc.LocalAddr()
	go func() {
		ns.server.ActivateAndServe()
		pc.Close()
	}()
	return nil
}

func (ns *nameserver) lookup(q string) []dns.RR {
	ns.Lock()
	defer ns.Unlock()
	return ns.rr[q]
}

func (ns *nameserver) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(req)
	q := msg.Question[0].Name // question, e.g. "miek.nl."

	rr := ns.lookup(q)
	if len(rr) == 0 { // no entry
		w.Close()
		return
	}

	msg.Answer = make([]dns.RR, 1)
	msg.Answer[0] = rr[0]
	w.WriteMsg(msg)
}

func (ns *nameserver) Scrub() {
	ns.Lock()
	defer ns.Unlock()
	if ns.server == nil {
		return
	}
	ns.server.Shutdown()
	ns.server = nil
}

func (ns *nameserver) X() circuit.X {
	return circuit.Ref(XNameserver{ns})
}

func (ns *nameserver) Set(rr string) error {
	ss, err := dns.NewRR(rr)
	if err != nil {
		return err
	}
	ns.Lock()
	defer ns.Unlock()
	ns.rr[ss.Header().Name] = append(ns.rr[ss.Header().Name], ss)
	return nil
}

func (ns *nameserver) Unset(name string) {
	ns.Lock()
	defer ns.Unlock()
	delete(ns.rr, name)
}

func (ns *nameserver) Peek() Stat {
	ns.Lock()
	defer ns.Unlock()
	var stat Stat
	stat.Address = ns.addr.String()
	for name, rr := range ns.rr {
		var ss []string
		for _, record := range rr {
			ss = append(ss, record.String())
		}
		stat.Records[name] = ss
	}
	return stat
}
