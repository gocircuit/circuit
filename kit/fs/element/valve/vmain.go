// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"encoding/json"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

func MakeValve() *Valve {
	return &Valve{ErrorFile: file.NewErrorFile()}
}

const MessageCap = 32e3 // 32K

type Valve struct {
	*file.ErrorFile
	sync.Mutex
	stat  Stat
	send  *senderValve
	recv  *receiverValve
}

type onceSender chan<- interface{}
type onceReceiver <-chan interface{}

// ctrl

func (c *Valve) SetCap(n int) error {
	c.ErrorFile.Clear()
	c.Lock()
	defer c.Unlock()
	if c.send != nil {
		c.ErrorFile.Set("capacity already set")
		return rh.ErrClash
	}
	if n < 0 {
		c.ErrorFile.Set("negative capacity")
		return rh.ErrPerm
	}
	c.stat.Cap = n
	c.stat.Opened = true
	ch := make(chan onceReceiver, n)
	c.send, c.recv = newSenderValve(c.ErrorFile, ch), newReceiverValve(c.ErrorFile, ch)
	return nil
}

func (c *Valve) GetCap() int {
	c.Lock()
	defer c.Unlock()
	return c.stat.Cap
}

// send-side

func (c *Valve) sendvalve() (*senderValve, error) {
	c.ErrorFile.Clear()
	c.Lock()
	defer c.Unlock()
	if c.send == nil {
		c.ErrorFile.Set("capacity not set")
		return nil, rh.ErrClash
	}
	return c.send, nil
}

func (c *Valve) Close() error {
	x, err := c.sendvalve()
	if err != nil {
		return err
	}
	c.Lock()
	c.stat.Closed = true
	c.Unlock()
	go func() {
		x.Close()
	}()
	return nil
}

func (c *Valve) incSend(err error) error {
	if err != nil {
		return err
	}
	c.Lock()
	defer c.Unlock()
	c.stat.NumSend++
	return nil
}

func (c *Valve) Send(v interface{}, intr rh.Intr) error {
	x, err := c.sendvalve()
	if err != nil {
		return err
	}
	return c.incSend(x.Send(v, intr))
}

func (c *Valve) TrySend(v interface{}) error {
	x, err := c.sendvalve()
	if err != nil {
		return err
	}
	return c.incSend(x.TrySend(v))
}

func (c *Valve) WaitSend(intr rh.Intr) error {
	x, err := c.sendvalve()
	if err != nil {
		return err
	}
	return x.WaitSend(intr)
}

// recv-side

func (c *Valve) recvvalve() (*receiverValve, error) {
	c.ErrorFile.Clear()
	c.Lock()
	defer c.Unlock()
	if c.recv == nil {
		c.ErrorFile.Set("capacity not set")
		return nil, rh.ErrClash
	}
	return c.recv, nil
}

func (c *Valve) incRecv(v interface{}, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}
	c.Lock()
	defer c.Unlock()
	c.stat.NumRecv++
	return v, nil
}

func (c *Valve) Recv(intr rh.Intr) (interface{}, error) {
	x, err := c.recvvalve()
	if err != nil {
		return nil, err
	}
	return c.incRecv(x.Recv(intr))
}

func (c *Valve) TryRecv() (interface{}, error) {
	x, err := c.recvvalve()
	if err != nil {
		return nil, err
	}
	return c.incRecv(x.TryRecv())
}

func (c *Valve) WaitRecv(intr rh.Intr) error {
	x, err := c.recvvalve()
	if err != nil {
		return err
	}
	return x.WaitRecv(intr)
}

// stat

type Stat struct {
	Cap     int  `json:"cap"`
	Opened  bool `json:"opened"`
	Closed  bool `json:"closed"`
	NumSend int  `json:"numsend"`
	NumRecv int  `json:"numrecv"`
}

func (s *Stat) String() string {
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (c *Valve) GetStat() *Stat {
	c.Lock()
	defer c.Unlock()
	var s Stat = c.stat
	return &s
}
