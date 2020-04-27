// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/hoijui/circuit/use/circuit"
)

type Valve interface {
	Send() (io.WriteCloser, error)
	IsDone() bool
	Scrub()
	Close() error
	Recv() (io.ReadCloser, error)
	Cap() int
	Stat() Stat
	X() circuit.X
}

// valve
type valve struct {
	send struct {
		abr <-chan struct{} // abort when closed
		sync.Mutex
		tun chan<- interface{}
	}
	recv struct {
		abr <-chan struct{} // abort when closed
		tun <-chan interface{}
	}
	ctrl struct {
		sync.Mutex
		abr  chan<- struct{}
		stat Stat
	}
}

type Stat struct {
	Cap     int  `json:"cap"`
	Opened  bool `json:"opened"`
	Closed  bool `json:"closed"`
	Aborted bool `json:"aborted"`
	NumSend int  `json:"numsend"`
	NumRecv int  `json:"numrecv"`
}

// Sender-receiver pipe capacity (once matched)
const MessageCap = 32e3 // 32K

func (s *Stat) String() string {
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func MakeValve(n int) Valve {
	v := &valve{}
	tun, abr := make(chan interface{}, n), make(chan struct{})
	v.send.tun, v.recv.tun = tun, tun
	v.ctrl.abr, v.send.abr, v.recv.abr = abr, abr, abr
	v.ctrl.stat.Opened, v.ctrl.stat.Cap = true, n
	return v
}

func (v *valve) X() circuit.X {
	return circuit.Ref(XValve{v})
}

func (v *valve) incSend() {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	v.ctrl.stat.NumSend++
}

func (v *valve) incRecv() {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	v.ctrl.stat.NumRecv++
}

// Cap returns the capacity of the valve and whether it was set.
func (v *valve) Cap() int {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	if v.ctrl.stat.Opened {
		return v.ctrl.stat.Cap
	}
	return -1
}

func (v *valve) Stat() Stat {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	return v.ctrl.stat
}
