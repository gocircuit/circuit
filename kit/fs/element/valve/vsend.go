// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"sync"

	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/interruptible"
)

// Close closes the channel
func (v *Valve) Close() error {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	if v.ctrl.stat.Closed {
		return rh.ErrGone
	}
	v.ctrl.stat.Closed = true
	close(v.ctrl.abr)
	go func() {
		v.send.Lock(nil) // Lock send system
		close(v.send.tun)
	}()
	return nil
}

func (v *Valve) Send(intr rh.Intr) (iw interruptible.Writer, err error) {
	// Lock send system
	println("valve.Send")
	u := v.send.Lock(intr)
	if u == nil {
		return nil, rh.ErrIntr
	}

	// Is there an abandoned gate to the receiver?
	if v.send.gate != nil {
		g := v.send.gate
		v.send.gate = nil
		return newValveWriter(v, u, g), nil
	}

	// Otherwise, make a new gate and push it to the receiver
	g := make(chan interruptible.Reader, 1)
	select {
	case v.send.tun <- g:
		return newValveWriter(v, u, g), nil
	case <-v.send.abr:
		u.Unlock()
		return nil, rh.ErrGone
	case <-intr:
		u.Unlock()
		return nil, rh.ErrIntr
	}
}

func (v *Valve) TrySend() (iw interruptible.Writer, err error) {
	// Lock send system
	u := v.send.TryLock()
	if u == nil {
		return nil, rh.ErrBusy
	}

	// Is there an abandoned gate to the receiver?
	if v.send.gate != nil {
		g := v.send.gate
		v.send.gate = nil
		return newValveWriter(v, u, g), nil
	}

	// Otherwise, make a new gate and push it to the receiver
	g := make(chan interruptible.Reader, 1)
	select {
	case v.send.tun <- g:
		return newValveWriter(v, u, g), nil
	case <-v.send.abr:
		u.Unlock()
		return nil, rh.ErrGone
	default:
		u.Unlock()
		return nil, rh.ErrBusy
	}
}

// valveWriter is an interruptible.Writer.
type valveWriter struct {
	valve *Valve
	u *interruptible.Unlocker
	s *Snatcher
	sync.Mutex
	g chan<- interruptible.Reader
	w interruptible.Writer
}

// ValveWriter …
// gate is the channel where the valve writer will send the reading end, if it commits to the session.
func newValveWriter(
	valve *Valve, 
	unlocker *interruptible.Unlocker,
	gate chan<- interruptible.Reader,
) (vw interruptible.Writer) {
	return &valveWriter{
		valve: valve,
		u: unlocker,
		s: NewSnatcher(),
		g: gate,
		w: brokenWriter{rh.ErrGone},
	}
}

const (
	committing = iota
	abandoning
)

func (vw *valveWriter) commit() interruptible.Writer {
	vw.Lock()
	defer vw.Unlock()
	if vw.s.Snatch(committing) == FirstSnatch {
		defer vw.u.Unlock()
		r, w := interruptible.Pipe()
		vw.g <- r
		vw.g = nil
		vw.w = w
	}
	return vw.w
}

func (vw *valveWriter) Write(p []byte) (int, error) {
	panic("not used")
}

func (vw *valveWriter) WriteIntr(p []byte, intr rh.Intr) (n int, err error) {
	return vw.commit().WriteIntr(p, intr)
}

func (vw *valveWriter) abandon() {
	vw.Lock()
	defer vw.Unlock()
	if vw.s.Snatch(abandoning) == FirstSnatch {
		defer vw.u.Unlock()
		vw.valve.send.gate, vw.g = vw.g, nil
	}
}

func (vw *valveWriter) Close() error {
	vw.abandon()
	return vw.commit().Close()
}

// brokenWriter is an interruptible.Writer which always fails in error.
type brokenWriter struct {
	error
}

func (broken brokenWriter) Write([]byte) (int, error) {
	return 0, broken.error
}

func (broken brokenWriter) WriteIntr(p []byte, intr rh.Intr) (n int, err error) {
	return 0, broken.error
}

func (broken brokenWriter) Close() error {
	return broken.error
}
