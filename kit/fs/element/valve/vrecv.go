// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"sync"

	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/interruptible"
)

type receiverValve struct {
	recv <-chan onceReceiver
	carrier struct {
		lk1 interruptible.Mutex
		lk2 sync.Mutex
		ch  onceReceiver
	}
	*file.ErrorFile
}

func newReceiverValve(errfile *file.ErrorFile, recv <-chan onceReceiver) *receiverValve {
	return &receiverValve{
		recv:      recv,
		ErrorFile: errfile,
	}
}

func (c *receiverValve) Recv(intr rh.Intr) (interface{}, error) {
	println("recv begin")
	defer println("recv end")

	c.ErrorFile.Clear()
	if t, ok := c.quickRecv(); ok {
		select {
		case v := <-t:
			return v, nil
		case <-intr:
			c.ErrorFile.Set("receive interrupted")
			return nil, rh.ErrIntr
		}
	}
	//
	select {
	case t, ok := <-c.recv:
		if !ok {
			return nil, rh.ErrGone
		}
		select {
		case v := <-t:
			return v, nil
		case <-intr:
			c.ErrorFile.Set("receive interrupted")
			return nil, rh.ErrIntr
		}
	case <-intr:
		c.ErrorFile.Set("receive interrupted")
		return nil, rh.ErrIntr
	}
}

func (c *receiverValve) TryRecv() (interface{}, error) {
	c.ErrorFile.Clear()
	if t, ok := c.quickRecv(); ok {
		select {
		case v := <-t:
			return v, nil
		default:
			// lose the message;
			// happens only when TrySend and TryRecv match each other;
			// this is a user programming mistake.
			c.ErrorFile.Set("try send matched try recv (user programming error); message lost")
			return nil, rh.ErrClash
		}
	}
	//
	select {
	case t, ok := <-c.recv:
		if !ok {
			return nil, rh.ErrGone
		}
		select {
		case v := <-t:
			return v, nil
		default:
			// lose the message;
			// happens only when TrySend and TryRecv match each other;
			// this is a user programming mistake.
			c.ErrorFile.Set("try send matched try recv (user programming error); message lost")
			return nil, rh.ErrClash
		}
	default:
		c.ErrorFile.Set("receive would block")
		return nil, rh.ErrBusy
	}
}

func (c *receiverValve) quickRecv() (onceReceiver, bool) {
	c.carrier.lk2.Lock()
	defer c.carrier.lk2.Unlock()
	if c.carrier.ch == nil {
		return nil, false
	}
	defer func() {
		c.carrier.ch = nil
	}()
	return c.carrier.ch, true
}

func (c *receiverValve) canQuickRecv() bool {
	c.carrier.lk2.Lock()
	defer c.carrier.lk2.Unlock()
	return c.carrier.ch != nil
}

func (c *receiverValve) WaitRecv(intr rh.Intr) error {
	u := c.carrier.lk1.Lock(intr)
	if u == nil {
		return rh.ErrIntr
	}
	defer u.Unlock()
	//
	c.ErrorFile.Clear()
	if c.canQuickRecv() {
		return nil
	}
	select {
	case t, ok := <-c.recv:
		if !ok {
			return rh.ErrGone
		}
		c.carrier.ch = t
		println("plant carrier")
		return nil
	case <-intr:
		c.ErrorFile.Set("send interrupted")
		return rh.ErrIntr
	}
}
