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

func (v *Valve) Recv(intr rh.Intr) (ir interruptible.Reader, err error) {
	// Lock send system
	u := v.recv.Lock(intr)
	if u == nil {
		return nil, rh.ErrIntr
	}

	// Is there an abandoned gate from the sender?
	if v.recv.gate != nil {
		g := v.recv.gate
		v.recv.gate = nil
		return newValveReader(v, u, g), nil
	}

	// Otherwise, pull a gate from the sender
	select {
	case g := <-v.recv.tun:
		return newValveReader(v, u, g), nil
	case <-v.recv.abr:
		u.Unlock()
		return nil, rh.ErrGone
	case <-intr:
		u.Unlock()
		return nil, rh.ErrIntr
	}
}

func (v *Valve) TryRecv() (ir interruptible.Reader, err error) {
	// Lock send system
	u := v.recv.TryLock()
	if u == nil {
		return nil, rh.ErrBusy
	}

	// Is there an abandoned gate from the sender?
	if v.recv.gate != nil {
		g := v.recv.gate
		v.recv.gate = nil
		return newValveReader(v, u, g), nil
	}

	// Otherwise, pull a gate from the sender
	select {
	case g, ok := <-v.recv.tun:
		if !ok {
			u.Unlock()
			return nil, rh.ErrEOF
		}
		return newValveReader(v, u, g), nil
	case <-v.recv.abr:
		u.Unlock()
		return nil, rh.ErrGone
	default:
		u.Unlock()
		return nil, rh.ErrBusy
	}
}

// valveReader is an interruptible.Reader.
type valveReader struct {
	valve *Valve
	u *interruptible.Unlocker
	s *Snatcher
	sync.Mutex
	g <-chan interruptible.Reader
	r interruptible.Reader
}

// ValveReader …
// gate is the channel where the valve writer will send the reading end, if it commits to the session.
func newValveReader(
	valve *Valve, 
	unlocker *interruptible.Unlocker,
	gate <-chan interruptible.Reader,
) (vr interruptible.Reader) {
	return &valveReader{
		valve: valve,
		u: unlocker,
		s: NewSnatcher(),
		g: gate,
		r: brokenReader{rh.ErrGone},
	}
}

func (vr *valveReader) commit(intr rh.Intr) interruptible.Reader {
	vr.Lock()
	defer vr.Unlock()
	if vr.s.Snatch(committing) == FirstSnatch {
		select {
		case vr.r = <-vr.g:
			vr.g = nil
			vr.u.Unlock()
		case <-intr:
			vr.r = brokenReader{rh.ErrIntr}
			vr.s = NewSnatcher()
		}
	}
	return vr.r
}

func (vr *valveReader) Read(p []byte) (int, error) {
	panic("not used")
}

func (vr *valveReader) ReadIntr(p []byte, intr rh.Intr) (n int, err error) {
	return vr.commit(intr).ReadIntr(p, intr)
}

func (vr *valveReader) abandon() {
	vr.Lock()
	defer vr.Unlock()
	if vr.s.Snatch(abandoning) == FirstSnatch {
		defer vr.u.Unlock()
		vr.valve.recv.gate, vr.g = vr.g, nil
	}
}

func (vr *valveReader) Close() error {
	vr.abandon()
	return vr.commit(nil).Close()
}

// brokenReader is an interruptible.Reader which always fails in error.
type brokenReader struct {
	error
}

func (broken brokenReader) Read([]byte) (int, error) {
	return 0, broken.error
}

func (broken brokenReader) ReadIntr(p []byte, intr rh.Intr) (n int, err error) {
	return 0, broken.error
}

func (broken brokenReader) Close() error {
	return broken.error
}
