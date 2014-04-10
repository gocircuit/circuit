// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"fmt"
	"net"
	"time"

	"github.com/gocircuit/circuit/kit/tele"
	"github.com/gocircuit/circuit/kit/tele/blend"
	"github.com/gocircuit/circuit/kit/tele/tcp"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

// Server
type Server struct {
	frame   trace.Frame
	tele    *blend.Transport
	inAddr  string
	outAddr string
}

func NewServer(inAddr, outAddr string) (*Server, error) {
	t := tele.NewStructOverTCP()
	l := t.Listen(tcp.Addr(outAddr))
	if outAddr == "" {
		outAddr = l.Addr().String()
		fmt.Println(outAddr)
	}
	srv := &Server{frame: trace.NewFrame("tele", "server"), tele: t, inAddr: inAddr, outAddr: outAddr}
	srv.frame.Bind(srv)
	go srv.loop(l)
	return srv, nil
}

func (srv *Server) loop(l *blend.Listener) {
	for {
		session := l.AcceptSession()
		go func() {
			for {
				outConn := session.Accept()
				// Read the first empty chunk from the connection
				if _, err := outConn.Read(); err != nil {
					srv.frame.Printf("first read (%s)", err)
					outConn.Close()
					continue
				}
				// Dial user server
				inConn, err := net.Dial("tcp", srv.inAddr)
				if err != nil {
					outConn.Close()
					srv.frame.Printf("server dial tcp address %s (%s)", srv.inAddr, err)
					time.Sleep(time.Second) // Prevents DoS when local TCP server is down temporarily
					continue
				}
				Proxy(inConn, outConn)
			}
		}()
	}
}
