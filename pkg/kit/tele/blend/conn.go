// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package blend

import (
	"io"
	"net"
	"sync"
)

/*

	 accept/dial            Conn
	+-----------+        +--------+
	|  pollread |------->| prompt |
	|    ···    |        |  ···   |
	|   write   |<-------| Write  |
	|   scrub   |<-------| Close  |
	+-----------+        +--------+

	Invariant "close-and-prompt": The connection object should not be
	registered with the AcceptConn as open, after the user has called Close,
	the AcceptConn has invoked Conn.prompt with an error.

	SOURCES OF CLOSURE:

	(nil,err) -----> prompt
	                  ···
	                 Close <--- USER
                      ···
	  write |<-----| Write
	        |-err->|

*/

type Conn struct {
	connID   ConnID
	ssn      ssn
	p__      sync.Mutex // send-side of prompt channel
	pch      chan *readReturn
	peof     bool // Prompt-side closure
	w__      sync.Mutex
	nwritten SeqNo // Number of writes
	weof     bool  // Write-side closure
}

type ssn interface {
	write(*Msg) error
	writeAbort(ConnID, error) error
	scrub(ConnID)
	RemoteAddr() net.Addr
}

type readReturn struct {
	Payload interface{}
	Err     error
}

func newConn(connID ConnID, ssn ssn) *Conn {
	return &Conn{connID: connID, ssn: ssn, pch: make(chan *readReturn, 3)}
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.ssn.RemoteAddr()
}

// Read reads the next chunk of bytes.
func (c *Conn) Read() (interface{}, error) {
	rr, ok := <-c.pch
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}
	return rr.Payload, rr.Err
}

func (c *Conn) prompt(payload interface{}, err error) {
	c.p__.Lock()
	defer c.p__.Unlock()
	if c.peof {
		return
	}
	c.pch <- &readReturn{Payload: payload, Err: err}
	if err != nil {
		close(c.pch)
		c.peof = true
	}
}

func (c *Conn) promptClose() {
	c.p__.Lock()
	defer c.p__.Unlock()
	if c.peof {
		return
	}
	close(c.pch)
	c.peof = true
}

// Write writes the chunk to the connection.
func (c *Conn) Write(v interface{}) error {
	c.w__.Lock()
	defer c.w__.Unlock()
	if c.weof {
		panic("writin after close")
	}
	c.nwritten++
	msg := &Msg{
		ConnID: c.connID,
		Demux: &PayloadMsg{
			SeqNo:   c.nwritten - 1,
			Payload: v,
		},
	}
	return c.ssn.write(msg)
}

// Close closes the connection. It is synchronized with Write and will not interrupt a concurring write.
func (c *Conn) Close() error {
	c.ssn.scrub(c.connID) // Scrub outside of w__ lock
	//
	c.w__.Lock()
	if c.weof {
		c.w__.Unlock()
		return io.ErrUnexpectedEOF
	}
	c.weof = true
	c.w__.Unlock()
	//
	c.promptClose()
	return nil
}

// Abort closes the connection ...
func (c *Conn) Abort(reason error) {
	if reason == nil {
		panic("x")
	}
	c.ssn.scrub(c.connID) // Scrub outside of w__ lock
	//
	c.w__.Lock()
	if c.weof {
		c.w__.Unlock()
		return
	}
	c.weof = true
	c.ssn.writeAbort(c.connID, reason)
	c.w__.Unlock()
	//
	c.promptClose()
	return
}
