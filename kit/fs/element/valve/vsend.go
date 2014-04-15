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

type senderValve struct {
	send struct {
		interruptible.Mutex
		send chan<- onceReceiver
	}
	carrier struct {
		lk1 interruptible.Mutex
		lk2 sync.Mutex
		ch  onceSender
	}
	*file.ErrorFile
}

func newSenderValve(errfile *file.ErrorFile, send chan<- onceReceiver) *senderValve {
	r := &senderValve{
		ErrorFile: errfile,
	}
	r.send.send = send
	return r
}

func (c *senderValve) Send(v interface{}, intr rh.Intr) error {
	c.ErrorFile.Clear()
	if c.quickSend(v) {
		return nil
	}
	//
	u := c.send.Lock(intr)
	if u == nil {
		return rh.ErrIntr
	}
	defer u.Unlock()
	//
	if c.send.send == nil {
		return rh.ErrGone
	}
	//
	t := make(chan interface{}, 1)
	t <- v
	select {
	case c.send.send <- t:
		return nil
	case <-intr:
		c.ErrorFile.Set("send interrupted")
		return rh.ErrIntr
	}
}

func (c *senderValve) TrySend(v interface{}) error {
	c.ErrorFile.Clear()
	if c.quickSend(v) {
		return nil
	}
	//
	u := c.send.TryLock()
	if u == nil {
		return rh.ErrIntr
	}
	defer u.Unlock()
	//
	if c.send.send == nil {
		return rh.ErrGone
	}
	//
	t := make(chan interface{}, 1)
	t <- v
	select {
	case c.send.send <- t:
		return nil
	default:
		c.ErrorFile.Set("send would block")
		return rh.ErrBusy
	}
}

func (c *senderValve) WaitSend(intr rh.Intr) error {
	u1 := c.carrier.lk1.Lock(intr)
	if u1 == nil {
		return rh.ErrIntr
	}
	defer u1.Unlock()
	//
	c.ErrorFile.Clear()
	if c.canQuickSend() {
		return nil
	}
	//
	u2 := c.send.Lock(intr)
	if u2 == nil {
		return rh.ErrIntr
	}
	defer u2.Unlock()
	//
	if c.send.send == nil {
		return rh.ErrGone
	}
	//
	t := make(chan interface{}, 1)
	select {
	case c.send.send <- t:
		c.carrier.ch = onceSender(t)
		return nil
	case <-intr:
		c.ErrorFile.Set("send interrupted")
		return rh.ErrIntr
	}
}

func (c *senderValve) Close() error {
	u := c.send.Lock(nil)
	defer u.Unlock()
	//
	if c.send.send == nil {
		return rh.ErrGone
	}
	//
	close(c.send.send)
	c.send.send = nil
	return nil
}

func (c *senderValve) quickSend(v interface{}) bool {
	c.carrier.lk2.Lock()
	defer c.carrier.lk2.Unlock()
	if c.carrier.ch == nil {
		return false
	}
	c.carrier.ch <- v
	c.carrier.ch = nil
	return true
}

func (c *senderValve) canQuickSend() bool {
	c.carrier.lk2.Lock()
	defer c.carrier.lk2.Unlock()
	return c.carrier.ch != nil
}
