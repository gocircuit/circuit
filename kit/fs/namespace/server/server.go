// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package server

import (
	"github.com/gocircuit/circuit/kit/fs/rh"
)

func NewServer(slash rh.FID) rh.Server {
	return &Server{slash}
}

type Server struct {
	slash rh.FID
}

func (s *Server) SignIn(user, dir string) (rh.Session, error) {
	return NewSession(s.slash), nil
}

func (s *Server) String() string {
	return "generic server"
}


func NewSession(slash rh.FID) rh.Session {
	return &Session{slash}
}

type Session struct {
	slash rh.FID
}

func (s *Session) Walk(name []string) (rh.FID, error) {
	return s.slash.Walk(name)
}

func (s *Session) SignOut() {}

func (s *Session) String() string {
	return "generic session"
}

