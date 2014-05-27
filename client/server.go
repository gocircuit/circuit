// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"io"
	"time"

	"github.com/gocircuit/circuit/element/srv"
)

// ServerStat encloses subscription state information.
type ServerStat struct {
	Addr string
	Joined time.Time
}

func srvStat(s srv.Stat) ServerStat {
	return ServerStat{
		Addr: s.Addr,
		Joined: s.Joined,
	}
}

// Server…
// All methods panic if the hosting circuit server dies.
type Server interface {
	Profile(string) (io.Reader, error)
	Peek() ServerStat
}

type ysrvSrv struct {
	srv.YServer
}

func (y ysrvSrv) Peek() ServerStat {
	return srvStat(y.YServer.Peek())
}
