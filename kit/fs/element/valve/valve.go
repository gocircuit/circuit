// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/interruptible"
)

func MakeValve() *Valve {
	v := &Valve{ErrorFile: file.NewErrorFile()}
	l := &loop{}

	send := make(chan interface{})
	v.send.send, l.send = send, send
	recv := make(chan chan<- interface{})
	v.recv.recv, l.recv = recv, recv

	sgone := make(chan struct{})
	v.send.gone, l.sgone = sgone, sgone
	rgone := make(chan struct{})
	v.recv.gone, l.rgone = rgone, rgone

	req, resp := make(chan interface{}), make(chan interface{})
	v.ctrl.req, l.req = req, req
	v.ctrl.resp, l.resp = resp, resp

	go l.main()
	return v
}

type Valve struct {
	send struct {
		interruptible.Mutex
		send chan<- interface{}
		gone <-chan struct{} // closure notification for Send
	}
	recv struct {
		interruptible.Mutex // sync competing access with option to interrupt
		recv chan<- chan<- interface{} // Recv sends a vessel for a message
		gone <-chan struct{}
	}
	ctrl struct {
		sync.Mutex // control ops are instantaneous
		req   chan<- interface{}
		resp  <-chan interface{}
	}
	*file.ErrorFile
}

var (
	ErrIntr = errors.New("intr")
	ErrGone = errors.New("gone")
)

func RhError(e error) error {
	switch e {
	case ErrIntr:
		return rh.ErrIntr
	case ErrGone:
		return rh.ErrGone
	case nil:
		return nil
	}
	panic(0)
}

type Stat struct {
	NSend  int `json:"numsend"` // total number of messages sent
	NRecv  int `json:"numrecv"` // total number of messages received
}

func (s *Stat) String() string {
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (c *Valve) Send(v interface{}, intr rh.Intr) error {
	c.ErrorFile.Clear()
	u := c.send.Mutex.Lock(intr)
	if u == nil {
		c.ErrorFile.Set("send interrupted")
		return ErrIntr
	}
	defer u.Unlock()
	select {
	case c.send.send <- v:
		return nil
	case <-c.send.gone:
		c.ErrorFile.Set("send on closed channel")
		return ErrGone
	case <-intr:
		c.ErrorFile.Set("send interrupted")
		return ErrIntr
	}
}

func (c *Valve) Recv(intr rh.Intr) (interface{}, error) {
	c.ErrorFile.Clear()
	u := c.recv.Mutex.Lock(intr)
	if u == nil {
		c.ErrorFile.Set("receive interrupted")
		return nil, ErrIntr
	}
	defer u.Unlock()
	var t = make(chan interface{})
	select {
	case c.recv.recv <- (chan<- interface{})(t):
		select {
		case v, ok := <-t:
			if !ok {
				panic(0)
			}
			return v, nil
		case <-c.recv.gone:
			c.ErrorFile.Set("receive on closed channel")
			return nil, ErrGone
		case <-intr:
			c.ErrorFile.Set("receive interrupted")
			return nil, ErrIntr
		}
	case <-c.recv.gone:
		c.ErrorFile.Set("receive on closed channel")
		return nil, ErrGone
	case <-intr:
		c.ErrorFile.Set("receive interrupted")
		return nil, ErrIntr
	}
}

func (c *Valve) Close() {
	c.ErrorFile.Clear()
	c.ctrl.Mutex.Lock()
	defer c.ctrl.Mutex.Unlock()
	if c.ctrl.req == nil {
		c.ErrorFile.Set("channel already closed")
		return
	}
	close(c.ctrl.req)
	c.ctrl.req = nil
}

func (c *Valve) IsClosed() bool {
	c.ctrl.Mutex.Lock()
	defer c.ctrl.Mutex.Unlock()
	return c.ctrl.req == nil
}

func (c *Valve) SetCap(n int) {
	c.ctrl.Mutex.Lock()
	defer c.ctrl.Mutex.Unlock()
	if c.ctrl.req == nil {
		return
	}
	c.ctrl.req <- reqSetCap{n}
}

func (c *Valve) GetCap() int {
	c.ctrl.Mutex.Lock()
	defer c.ctrl.Mutex.Unlock()
	if c.ctrl.req == nil {
		return -1
	}
	c.ctrl.req <- reqGetCap{}
	return (<-c.ctrl.resp).(int)
}

func (c *Valve) GetStat() *Stat {
	c.ctrl.Mutex.Lock()
	defer c.ctrl.Mutex.Unlock()
	if c.ctrl.req == nil {
		return nil
	}
	c.ctrl.req <- reqGetStat{}
	return (<-c.ctrl.resp).(*Stat)
}

func (c *Valve) WaitSend(intr rh.Intr) error {
	c.ErrorFile.Clear()
	c.ctrl.Mutex.Lock()
	if c.ctrl.req == nil {
		c.ctrl.Mutex.Unlock()
		c.ErrorFile.Set("waiting for send on a closed channel")
		return rh.ErrGone
	}
	unblock := make(chan struct{})
	c.ctrl.req <- reqWaitSend{unblock}
	c.ctrl.Mutex.Unlock()
	//
	select {
	case <-unblock:
		return nil
	case <-intr:
		c.ErrorFile.Set("waiting for send interrupted")
		return rh.ErrIntr
	}
}

func (c *Valve) WaitRecv(intr rh.Intr) error {
	c.ErrorFile.Clear()
	c.ctrl.Mutex.Lock()
	if c.ctrl.req == nil {
		c.ctrl.Mutex.Unlock()
		c.ErrorFile.Set("waiting for receive on a closed channel")
		return rh.ErrGone
	}
	unblock := make(chan struct{})
	c.ctrl.req <- reqWaitRecv{unblock}
	c.ctrl.Mutex.Unlock()
	//
	select {
	case <-unblock:
		return nil
	case <-intr:
		c.ErrorFile.Set("waiting for receive interrupted")
		return rh.ErrIntr
	}
}
