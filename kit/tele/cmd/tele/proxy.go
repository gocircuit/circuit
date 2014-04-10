// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"net"

	"github.com/gocircuit/circuit/kit/tele/blend"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

type proxy struct {
	frame  trace.Frame
	legacy net.Conn
	tele   *blend.Conn
}

func Proxy(legacy net.Conn, tele *blend.Conn) {
	p := &proxy{frame: trace.NewFrame("proxy"), legacy: legacy, tele: tele}
	p.frame.Bind(p)
	go p.legacy2tele()
	go p.tele2legacy()
}

const ReadBlockLen = 1e4

func (p *proxy) legacy2tele() {
	var (
		n   int
		err error
	)
	// We avoid buffer creation on each iteration since blend.Write copies the data before it returns.
	buf := make([]byte, ReadBlockLen)
	for {
		n, err = p.legacy.Read(buf)
		if n > 0 {
			if err := p.tele.Write(&cargo{Cargo: buf[:n]}); err != nil {
				p.frame.Printf("write (%s)", err)
				p.tele.Close()
				return
			}
			continue
		}
		if err == nil {
			panic("e")
		}
		p.frame.Printf("read (%s)", err)
		p.tele.Close()
		return
	}
}

func (p *proxy) tele2legacy() {
	var (
		err   error
		chunk interface{}
	)
	for {
		chunk, err = p.tele.Read()
		if err != nil {
			p.frame.Printf("read (%s)", err)
			p.legacy.Close()
			return
		}
		if _, err = p.legacy.Write(chunk.(*cargo).Cargo); err != nil {
			p.frame.Printf("write (%s)", err)
			p.legacy.Close()
			return
		}
	}
}
