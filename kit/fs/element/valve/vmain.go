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
	"github.com/gocircuit/circuit/kit/interruptible"
)

func MakeValve() *Valve {
	return &Valve{
		ErrorFile: file.NewErrorFile(),
	}
}

// Sender-receiver pipe capacity (once matched)
const MessageCap = 32e3 // 32K

// Valve
type Valve struct {
	send struct {
		abr <-chan struct{} // abort when closed
		interruptible.Mutex
		tun chan<- chan interruptible.Reader
		gate chan<- interruptible.Reader
	}
	recv struct {
		abr <-chan struct{} // abort when closed
		interruptible.Mutex
		tun <-chan chan interruptible.Reader
		gate <-chan interruptible.Reader
	}
	ErrorFile *file.ErrorFile
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

func (v *Valve) incSend() {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	v.ctrl.stat.NumSend++
}

func (v *Valve) incRecv() {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	v.ctrl.stat.NumRecv++
}

// SetCap initializes the valve and sets its buffer capacity to n.
func (v *Valve) SetCap(n int) error {
	// Lock the send system
	su := v.send.TryLock()
	if su == nil {
		v.ErrorFile.Set("concurring send attempt")
		return rh.ErrPerm
	}
	defer su.Unlock()

	// Lock the receive system
	ru := v.recv.TryLock()
	if ru == nil {
		v.ErrorFile.Set("concurring receive attempt")
		return rh.ErrPerm
	}
	defer ru.Unlock()

	// Lock the control system
	v.ctrl.Lock()
	defer v.ctrl.Unlock()

	// Validate argument and check we are not opening twice
	if v.ctrl.stat.Opened {
		v.ErrorFile.Set("capacity already set")
		return rh.ErrClash
	}
	if n < 0 {
		v.ErrorFile.Set("negative capacity")
		return rh.ErrPerm
	}

	// Initialize valve
	tun, abr := make(chan chan interruptible.Reader, n), make(chan struct{})
	v.send.tun, v.recv.tun = tun, tun                  // setup main channel
	v.ctrl.abr, v.send.abr, v.recv.abr = abr, abr, abr // setup abort channel
	v.ctrl.stat.Opened, v.ctrl.stat.Cap = true, n                // update state

	return nil
}

// GetCap returns the capacity of the valve and whether it was set.
func (v *Valve) GetCap() int {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	if v.ctrl.stat.Opened {
		return v.ctrl.stat.Cap
	}
	return -1
}

// GetStat …
func (v *Valve) GetStat() *Stat {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	var s Stat = v.ctrl.stat
	return &s
}
