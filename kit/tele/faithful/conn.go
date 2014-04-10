// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package faithful provides a lossless chunked connection over a lossy one.
package faithful

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/gocircuit/circuit/kit/tele/chain"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

type Conn struct {
	frame      trace.Frame
	sub        *chain.Conn
	ach        chan<- struct{}    // Abort channel
	rch        <-chan interface{} // []byte or error
	sync.Mutex                    // Linearizes Write requests
	bfr        *Buffer            // Buffer is self-synchronized // Write, Close
}

type control struct {
	Writer *chain.ConnWriter
	Msg    encoder
}

func NewConn(frame trace.Frame, under *chain.Conn) *Conn {

	// Only 1 needed in readLoop; get 3 just to be safe
	rch := make(chan interface{}, 3)
	// Capacity 1 (below) unblocks writeSync, invoked in readLoop right after a connection stitch
	// stitch is received, racing with a user write waiting on waitForLink.
	// In particular, if execUserWrite is waiting on waitForLink, it would prevent readLoop from
	// moving on to adopt the new connection.
	sch := make(chan *control, 1)
	// Abort channel
	ach := make(chan struct{})

	// User-facing Conn
	c := &Conn{
		frame: frame,
		sub:   under,
		rch:   rch,
		ach:   ach,
		bfr:   NewBuffer(frame.Refine("buffer"), MemoryCap),
	}
	c.frame.Bind(c)

	// readConn
	rc := &readConn{
		frame: frame.Refine("R∞"),
		sub:   c.sub,
		rch:   rch,
		sch:   sch,
		ach:   ach,
		bfr:   c.bfr,
	}
	rc.frame.Bind(rc)
	go rc.loop()

	// writeConn
	wc := &writeConn{
		frame: frame.Refine("W∞"),
		sub:   c.sub,
		bfr:   c.bfr,
		sch:   sch,
	}
	wc.frame.Bind(wc)
	go wc.loop()
	return c
}

func (c *Conn) Debug() interface{} {
	return c.sub.Debug()
}

const (
	MemoryCap      = 40
	AckFrequency   = 20
	LingerDuration = 30 * time.Second
)

// RemoteAddr returns the address of the remote endpoint on this faithful connection.
func (c *Conn) RemoteAddr() net.Addr {
	return c.sub.RemoteAddr()
}

// Read returns the next chunk received on the connection.
// chunk is non-nil if and only if err is nil.
func (c *Conn) Read() (chunk []byte, err error) {
	chunkOrErr, ok := <-c.rch
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}
	chunk, ok = chunkOrErr.([]byte)
	if ok {
		return chunk, nil
	}
	return nil, chunkOrErr.(error)
}

// Write blocks until the chunk is written to the connection.
// Write never returns an error, unless the connection is permanently broken.
func (c *Conn) Write(chunk []byte) (err error) {
	// Linearize user Writes & Closes
	c.Lock()
	defer c.Unlock()
	msg := &Chunk{chunk: chunk}
	return c.bfr.Write(msg)
}

// Close closes the connection semantically. The connection object will linger for a short while
// to ensure that the closure event is delivered to the remote endpoint.
func (c *Conn) Close() (err error) {
	c.Lock()
	defer c.Unlock()
	c.bfr.Close()
	go func() {
		c.frame.Println("linger starting")
		<-time.NewTimer(LingerDuration).C
		c.frame.Println("linger expired, carrier closed")
		c.sub.Kill()
		close(c.ach)
	}()
	return nil
}
