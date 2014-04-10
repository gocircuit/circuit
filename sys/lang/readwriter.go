// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"bytes"
	"encoding/gob"
	"io"
	"net"
	"sync"

	"github.com/gocircuit/circuit/use/n"
)

func NewBytesConn(addr string) n.Conn {
	var b bytes.Buffer
	return ReadWriterConn(stringAddr(addr), nopCloser{&b})
}

type nopCloser struct {
	io.ReadWriter
}

func (nc nopCloser) Close() error {
	return nil
}

type stringAddr string

func (a stringAddr) WorkerID() n.WorkerID {
	return ""
}

func (a stringAddr) NetAddr() net.Addr {
	return nil
}

func (a stringAddr) String() string {
	return string(a)
}

func (a stringAddr) FileName() string {
	return string(a)
}

// ReadWriterConn converts an io.ReadWriteCloser into a Conn
func ReadWriterConn(addr n.Addr, rwc io.ReadWriteCloser) n.Conn {
	return &readWriterConn{
		addr: addr,
		rwc:  rwc,
		enc:  gob.NewEncoder(rwc),
		dec:  gob.NewDecoder(rwc),
	}
}

type readWriterConn struct {
	addr n.Addr
	sync.Mutex
	rwc io.ReadWriteCloser
	enc *gob.Encoder
	dec *gob.Decoder
}

type blob struct {
	Cargo interface{}
}

func (conn *readWriterConn) Read() (interface{}, error) {
	conn.Lock()
	defer conn.Unlock()
	var b blob
	err := conn.dec.Decode(&b)
	if err != nil {
		return nil, err
	}
	return b.Cargo, nil
}

func (conn *readWriterConn) Write(cargo interface{}) error {
	conn.Lock()
	defer conn.Unlock()
	return conn.enc.Encode(&blob{cargo})
}

func (conn *readWriterConn) Close() error {
	conn.Lock()
	defer conn.Unlock()
	return conn.rwc.Close()
}

func (conn *readWriterConn) Abort(error) {
	conn.Lock()
	defer conn.Unlock()
	conn.rwc.Close()
}

func (conn *readWriterConn) Addr() n.Addr {
	return conn.addr
}
