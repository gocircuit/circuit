// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"github.com/hoijui/circuit/element/dns"
)

// NameserverStat encloses process state information.
type NameserverStat struct {

	// IP address of the nameserver
	Address string

	// Resource records resolved by this nameserver
	Records map[string][]string
}

func nameserverStat(s dns.Stat) NameserverStat {
	return NameserverStat{
		Address: s.Address,
		Records: s.Records,
	}
}

type Nameserver interface {

	Set(rr string) error

	Unset(name string)

	// Peek asynchronously returns the current state of the server.
	Peek() NameserverStat

	// Scrub shuts down the nameserver and removes its circuit element.
	Scrub()
}

type yNameserver struct {
	dns.YNameserver
}

func (y yNameserver) Peek() NameserverStat {
	return nameserverStat(y.YNameserver.Peek())
}
