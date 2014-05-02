// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"sync"
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
		v.send.Lock()
		defer v.send.Unlock()
		close(v.send.tun)
	}()
	return nil
}

// Send …
func (v *Valve) Send() (io.WriteCloser, err error) {
	v.send.Lock()
	defer v.send.Unlock()

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
		v.incSend()
		return newValveWriter(v, u, g), nil
	case <-v.send.abr:
		u.Unlock()
		return nil, rh.ErrGone
	case <-intr:
		u.Unlock()
		return nil, rh.ErrIntr
	}
}
