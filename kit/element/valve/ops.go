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
func (v *Valve) Send() (io.WriteCloser, error) {
	v.send.Lock()
	defer v.send.Unlock()
	r, w := interruptible.Pipe()
	select {
	case v.send.tun <- r:
		v.incSend()
		return w, nil
	case <-v.send.abr:
		return nil, errors.New("channel aborted")
	}
}

// Close closes the channel
func (v *Valve) Close() error {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	if v.ctrl.stat.Closed {
		return errors.New("channel already closed")
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

func (v *Valve) Recv() (io.ReadCloser, error) {
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
