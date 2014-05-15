// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"errors"
	"io"

	"github.com/gocircuit/circuit/kit/interruptible"
)

// Send …
// The returned WriteCloser must be closed at finalization.
func (v *valve) Send() (io.WriteCloser, error) {
	v.send.Lock()
	defer v.send.Unlock()
	if v.send.tun == nil {
		return nil, errors.New("channel closed")
	}
	r, w := interruptible.BufferPipe(32e3)
	select {
	case v.send.tun <- r:
		v.incSend()
		return w, nil
	case <-v.send.abr:
		return nil, errors.New("channel aborted")
	}
}

func (v *valve) IsDone() bool {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	return (v.ctrl.stat.Closed && v.ctrl.stat.NumSend == v.ctrl.stat.NumRecv) || v.ctrl.stat.Aborted
}

func (v *valve) Scrub() {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	if v.ctrl.stat.Aborted {
		return
	}
	close(v.ctrl.abr)
	v.ctrl.stat.Aborted = true
}

// Close closes the channel
func (v *valve) Close() error {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	if v.ctrl.stat.Closed {
		return errors.New("channel already closed")
	}
	v.ctrl.stat.Closed = true
	// The goroutine ensures close returns instantaneously, even if the
	// user erroneously races a send with it. In the latter case, the racing
	// sends will finish before closure takes place.
	go func() {
		v.send.Lock()
		defer v.send.Unlock()
		close(v.send.tun)
		v.send.tun = nil
	}()
	return nil
}

func (v *valve) Recv() (io.ReadCloser, error) {
	select {
	case g, ok := <-v.recv.tun:
		if !ok {
			return nil, errors.New("channel closed")
		}
		v.incRecv()
		return g.(interruptible.Reader), nil
	case <-v.recv.abr:
		return nil, errors.New("channel aborted")
	}
}
