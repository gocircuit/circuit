// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package server

import (
	"errors"
	"io"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hoijui/circuit/kit/interruptible"
	"github.com/hoijui/circuit/tissue"
	"github.com/hoijui/circuit/use/circuit"
	"github.com/hoijui/circuit/use/n"
)

type Server interface {
	Profile(string) (io.ReadCloser, error)
	Peek() Stat
	Rejoin(string) error // circuit address to join to
	Suicide()
	IsDone() bool
	Scrub()
	X() circuit.X
}

// server
type server struct {
	addr   string
	kin    *tissue.Kin
	joined time.Time
}

func New(kin *tissue.Kin) Server {
	return &server{
		addr:   kin.Avatar().X.Addr().String(),
		kin:    kin,
		joined: time.Now(),
	}
}

type Stat struct {
	Addr   string
	Joined time.Time
}

func (s *server) Rejoin(addr string) error {
	a, err := n.ParseAddr(addr)
	if err != nil {
		return err
	}
	return s.kin.ReJoin(a)
}

func (s *server) Suicide() {
	os.Exit(0)
}

func (s *server) Profile(name string) (io.ReadCloser, error) {
	p := pprof.Lookup(name)
	if p == nil {
		return nil, errors.New("no profile")
	}
	r, w := interruptible.Pipe()
	go func() {
		p.WriteTo(w, 1)
		w.Write([]byte("•••\n"))
		w.Close()
	}()
	return r, nil
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}

func (s *server) Peek() Stat {
	return Stat{
		Addr:   s.addr,
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
