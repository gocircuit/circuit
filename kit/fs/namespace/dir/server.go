// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package dir

import (
	"fmt"

	"github.com/gocircuit/circuit/kit/fs/rh"
)

// Server
type Server struct {
	slash rh.FID
}

func NewServer(slash rh.FID) *Server {
	return &Server{
		slash: slash,
	}
}

func (srv *Server) SignIn(user, _ string) (rh.Session, error) {
	return &Session{
		srv:   srv,
		user:  user,
		slash: srv.slash,
	}, nil
}

func (srv *Server) String() string {
	return fmt.Sprintf("rh/server(%s)", srv.slash.Q().String())
}

// Session
type Session struct {
	srv   *Server
	user  string
	slash rh.FID
}

func (ssn *Session) String() string {
	return fmt.Sprintf("rh/session(%s)", ssn.slash.Q().String())
}

func (ssn *Session) Walk(wname []string) (rh.FID, error) {
	return ssn.slash.Walk(wname)
}

func (ssn *Session) SignOut() {}
