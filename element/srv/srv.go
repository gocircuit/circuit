// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package srv

import (
	"errors"
	"io"
	"runtime/pprof"
	"time"

	"github.com/gocircuit/circuit/kit/interruptible"
	"github.com/gocircuit/circuit/use/circuit"
)

type Server interface {
	Profile(string) (io.Reader, error)
	Peek() Stat
	IsDone() bool
	Scrub()
	X() circuit.X
}

// server
type server struct {
	addr string
	joined time.Time
}

func New(addr string) Server {
	return &server{
		addr: addr,
		joined: time.Now(),
	}
}

type Stat struct {
	Addr string
	Joined time.Time
}

func (s *server) Profile(name string) (io.Reader, error) {
	p := pprof.Lookup(name)
	if p == nil {
		return nil, errors.New("no profile")
	}
	r, w := interruptible.BufferPipe(8e3)
	go p.WriteTo(w, 0)
	return r, nil
}

func (s *server) Peek() Stat {
	return Stat{
		Addr: s.addr,
		Joined: s.joined,
	}
}

func (s *server) IsDone() bool {
	return false
}

func (s *server) Scrub() {}

func (s *server) X() circuit.X {
	return circuit.Ref(XServer{s})
}
