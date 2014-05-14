// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package iomisc

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"

	"github.com/gocircuit/circuit/use/errors"
)

// ??
func ForwardClose(_ string, sink io.WriteCloser, source io.Reader, eof func()) {
	z := make(chan struct{}, 2)
	go func() {
		io.Copy(sink, source)
		sink.Close()
		z <- struct{}{}
	}()
	go func() {
		<-z
		if eof != nil {
			eof()
		}
	}()
}

// ??
func SniffClose(name string, sink io.WriteCloser, source io.Reader, eof func()) {
	sniffer := newSniffer(name)
	z1, z2 := make(chan struct{}), make(chan struct{})
	go func() {
		io.Copy(sniffer, source)
		sniffer.warn("SRC=>EVE⋅⋅SINK", "closing")
		sniffer.Close()
		close(z1)
	}()
	go func() {
		io.Copy(sink, sniffer)
		sniffer.warn("SRC⋅⋅EVE=>SINK", "closing")
		sink.Close()
		close(z2)
	}()
	go func() {
		<-z1
		sniffer.warn("SRC=>EVE⋅⋅SINK", "closed")
		<-z2
		sniffer.warn("SRC⋅⋅EVE=>SINK", "closed")
		if eof != nil {
			eof()
		}
	}()
}

//
type Sniffer struct {
	name string
	sync.Mutex
	closed   bool
	buf      bytes.Buffer
	unblock  chan struct{}
	nblocked int
}

func newSniffer(name string) *Sniffer {
	s := &Sniffer{
		name:    name,
		unblock: make(chan struct{}),
	}
	runtime.SetFinalizer(s, func(x *Sniffer) {
		x.warn("•", "Finalizer")
		x.Close()
	})
	return s
}

var errWait = errors.NewError("wait")

func (s *Sniffer) Printf(format string, arg ...interface{}) {
	log.Printf("▒ %s ▒ %s", s.name, fmt.Sprintf(format, arg...))
}

func (s *Sniffer) report(attr string, p []byte) {
	s.Printf("%s⟩\n“%s”", attr, string(p))
}

func (s *Sniffer) warn(attr, msg string) {
	s.Printf("%s⟫ “%s”", attr, msg)
}

func (s *Sniffer) Read(p []byte) (n int, err error) {
	for {
		if n, err = s.read(p); err != errWait {
			//s.report("R", p[:n])
			return
		}
		if _, ok := <-s.unblock; !ok {
			s.warn("R", "XEOF")
			return 0, io.ErrUnexpectedEOF
		}
	}
}

func (s *Sniffer) read(p []byte) (n int, err error) {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return 0, io.ErrUnexpectedEOF
	}
	n, _ = s.buf.Read(p)
	if n == 0 {
		s.nblocked++
		return 0, errWait
	}
	return
}

func (s *Sniffer) Write(p []byte) (n int, err error) {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return 0, io.ErrUnexpectedEOF
	}
	s.buf.Write(p)
	for ; s.nblocked > 0; s.nblocked-- {
		s.unblock <- struct{}{}
	}
	//s.report("W", p)
	return len(p), nil
}

func (s *Sniffer) Close() (err error) {
	s.Lock()
	defer s.Unlock()
	s.warn("C", "CLOSE")
	if s.closed {
		return io.ErrUnexpectedEOF
	}
	s.closed = true
	close(s.unblock)
	return nil
}
