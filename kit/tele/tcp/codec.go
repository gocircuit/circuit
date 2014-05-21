// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package tcp implements carrier transports over TCP.
package tcp

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
	"strings"

	"github.com/gocircuit/circuit/kit/tele/codec"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

// CodecTransport is a codec.Carrier over TCP.
var CodecTransport = codecTransport{trace.NewFrame("tcp")}

type codecTransport struct {
	trace.Frame
}

func (codecTransport) Listen(addr net.Addr) codec.CarrierListener {
	t := addr.String()
	if strings.Index(t, ":") < 0 {
		t = t + ":0"
	}
	l, err := net.Listen("tcp", t)
	if err != nil {
		return nil
	}
	return codecListener{l}
}

func (codecTransport) Dial(addr net.Addr) (codec.CarrierConn, error) {
	c, err := net.Dial("tcp", addr.String())
	if err != nil {
		return nil, err
	}
	return newCodecConn(trace.NewFrame("tcp", "dial"), c.(*net.TCPConn)), nil
}

type codecListener struct {
	net.Listener
}

func (l codecListener) Addr() net.Addr {
	return l.Listener.Addr()
}

func (l codecListener) Accept() (codec.CarrierConn) {
	c, err := l.Listener.Accept()
	if err != nil {
		log.Printf("error accepting tcp connection: %v", err)
		return nil
	}
	return newCodecConn(trace.NewFrame("tcp", "acpt"), c.(*net.TCPConn))
}

type codecConn struct {
	trace.Frame
	tcp *net.TCPConn
	r *bufio.Reader
}

func newCodecConn(f trace.Frame, c *net.TCPConn) *codecConn {
	if err := c.SetKeepAlive(true); err != nil {
		panic(err)
	}
	return &codecConn{f, c, bufio.NewReader(c)}
}

func (c *codecConn) RemoteAddr() net.Addr {
	return c.tcp.RemoteAddr()
}

func (c *codecConn) Read() (chunk []byte, err error) {
	k, err := binary.ReadUvarint(c.r)
	if err != nil {
		return nil, err
	}
	q := make([]byte, k)
	n, err := c.r.Read(q)
	return q[:n], err
}

func (c *codecConn) Write(chunk []byte) (err error) {
	q := make([]byte, len(chunk)+8)
	n := binary.PutUvarint(q, uint64(len(chunk)))
	m := copy(q[n:], chunk)
	_, err = c.tcp.Write(q[:n+m])
	return err
}

func (c *codecConn) Close() (err error) {
	return c.tcp.Close()
}
