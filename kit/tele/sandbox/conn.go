// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package sandbox

import (
	"bytes"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gocircuit/circuit/kit/tele/trace"
)

// Conn implements one net.Conn endpoint of a communication pipe.
//
//		+---+                 +---+
//		| R |<----------------| W |
//		|   |                 |   |
//		| W |---------------->| R |
//		+---+                 +---+
//		HalfConn              HalfConn
//
type HalfConn struct {
	frame   trace.Frame
	laddr   net.Addr
	raddr   net.Addr
	in, out *Sieve // Set in StartSieve
	//
	kill__ sync.Mutex
	kch    chan struct{}
	killed bool
	//
	recv__ sync.Mutex
	recv   <-chan interface{} // []byte or eof
	buf    bytes.Buffer
	//
	send__ sync.Mutex
	send   chan<- interface{} // []byte or eof
	fbck   <-chan feedback    // Feedback on the success of writes; nil if not supported
	dlvr   bool               // Back-channel information whether the last write was actually received
}

type eof struct{}

// NewHalfConn creates a new unbound sandbox connection endpoint.
func NewHalfConn(f trace.Frame) *HalfConn {
	hc := &HalfConn{frame: f, kch: make(chan struct{})}
	hc.frame.Bind(hc)
	return hc
}

type DebugInfo struct {
	In, Out *Sieve
}

func (p *HalfConn) Debug() interface{} {
	return &DebugInfo{In: p.in, Out: p.out}
}

// RecvFrom binds this endpoint to receive from (recv, addr).
func (p *HalfConn) RecvFrom(recv <-chan interface{}, addr net.Addr) {
	p.recv, p.raddr = recv, addr
}

// SendTo binds this endpoint to send to (send,addr) and consume feedback from fbck, if not nil.
func (p *HalfConn) SendTo(send chan<- interface{}, fbck <-chan feedback, addr net.Addr) {
	p.send, p.fbck, p.laddr = send, fbck, addr
}

// Read reads from the connection. Graceful closure on the remote side is returned as io.EOF.
func (p *HalfConn) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		panic("trivial buffer")
	}
	p.recv__.Lock()
	defer p.recv__.Unlock()

	// Read left-over data
	if n, _ = p.buf.Read(b); n > 0 {
		return n, nil
	}
	if p.buf.Len() != 0 {
		panic("u")
	}

	// Receive more data if capacity left in the user read buffer.
	// Close on this endpoint should not interrupt reading; reading has to be closed explicitly by the remote endpoint.
	t, ok := <-p.recv
	if !ok {
		return 0, io.ErrUnexpectedEOF
	}
	switch q := t.(type) {
	case []byte:
		if len(q) == 0 {
			panic("eh")
		}
		//p.frame.Printf("Read: #%v", q)
		// Write new data to internal read buffer
		if m, err := p.buf.Write(q); err != nil || m != len(q) {
			log.Printf("buf write %d/%d (%v)", m, len(q), err)
			panic("u/e")
		}
		// Read new data into user read buffer
		n, _ = p.buf.Read(b)
		return n, nil
	case eof:
		return 0, io.EOF
	}
	panic("eh")
}

// Delivered returns true if the last Write was delivered to the remote destination (i.e. it was not dropped).
func (p *HalfConn) Delivered() bool {
	p.send__.Lock()
	defer p.send__.Unlock()
	return p.dlvr
}

func (p *HalfConn) Write(b []byte) (n int, err error) {
	p.send__.Lock()
	defer p.send__.Unlock()
	if p.send == nil {
		return 0, io.ErrUnexpectedEOF
	}
	q := make([]byte, len(b))
	copy(q, b)
	select {
	case p.send <- q:
		n = len(q)
		// After a successful write to send, we are obligated to do one read from feedback.
		if p.fbck != nil {
			f := <-p.fbck
			p.dlvr = f.OK
			if f.EOF {
				err = io.ErrUnexpectedEOF
			}
		} else {
			p.dlvr = true
		}
	case <-p.kch:
		err = io.ErrUnexpectedEOF
	}
	return
}

// Close should be enqueued on the write path, so it doesn't interrupt pending writes
func (p *HalfConn) Close() error {
	if err := p.killWrite(); err != nil {
		return err
	}
	p.send__.Lock()
	defer p.send__.Unlock()
	if p.send == nil {
		return io.ErrUnexpectedEOF
	}
	// After a successful write to send, we are obligated to do one read from feedback.
	p.send <- eof{}
	if p.fbck != nil {
		<-p.fbck
	}
	close(p.send)
	p.send = nil
	return nil
}

func (p *HalfConn) killWrite() error {
	p.kill__.Lock()
	defer p.kill__.Unlock()
	if p.killed {
		return io.ErrUnexpectedEOF
	}
	close(p.kch)
	p.killed = true
	return nil
}

func (p *HalfConn) LocalAddr() net.Addr {
	return p.laddr
}

func (p *HalfConn) RemoteAddr() net.Addr {
	return p.raddr
}

func (p *HalfConn) SetDeadline(t time.Time) error {
	panic("n/s")
}

func (p *HalfConn) SetReadDeadline(t time.Time) error {
	panic("n/s")
}

func (p *HalfConn) SetWriteDeadline(t time.Time) error {
	panic("n/s")
}
