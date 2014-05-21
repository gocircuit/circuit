// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tele

import (
	// "log"
	"sync"
	"time"

	"github.com/gocircuit/circuit/kit/tele/blend"
	"github.com/gocircuit/circuit/use/errors"
	"github.com/gocircuit/circuit/use/n"
)

// Dialer
type Dialer struct {
	dialback n.Addr
	sub      *blend.Transport // Encloses *blend.Dialer
	sync.Mutex
	open map[n.WorkerID]*blend.DialSession // Open dial sessions
}

func newDialer(dialback n.Addr, sub *blend.Transport) *Dialer {
	return &Dialer{
		dialback: dialback,
		sub:      sub,
		open:     make(map[n.WorkerID]*blend.DialSession),
	}
}

func (d *Dialer) Dial(addr n.Addr) (conn n.Conn, err error) {
	d.Lock()
	defer d.Unlock()
	//
	workerID := addr.WorkerID()
	s, present := d.open[workerID]
	if !present {
		// Make new session to worker if one not present
		s, err = d.sub.DialSession(addr.(*Addr).TCP, func() {
			d.scrub(addr.WorkerID())
		})
		if err != nil {
			return nil, err
		}
		if err = d.auth(addr, s.Dial()); err != nil {
			s.Close()
			return nil, err
		}
		d.open[workerID] = s
		go d.watch(workerID, s) // Watch for idleness and close
	}
	return NewConn(s.Dial(), addr.(*Addr)), nil
}

// Idleness duration should be greater than the locus heartbeats over permanent cross-references
const IdleDuration = time.Second * 10

func (d *Dialer) watch(workerID n.WorkerID, s *blend.DialSession) {
	var ready bool
	for {
		time.Sleep(IdleDuration)
		if d.expire(workerID, s, &ready) {
			return
		}
	}
}

func (d *Dialer) expire(workerID n.WorkerID, s *blend.DialSession, ready *bool) (closed bool) {
	d.Lock()
	defer d.Unlock()
	//
	numconn, lastuse := s.NumConn()
	if numconn == 0 && time.Now().Sub(lastuse) > IdleDuration {
		if *ready {
			delete(d.open, workerID)
			// log.Printf("idle session %s expiring", s)
			s.Close()
			return true
		}
		*ready = true
	}
	return false
}

func (d *Dialer) scrub(workerID n.WorkerID) {
	d.Lock()
	defer d.Unlock()
	delete(d.open, workerID)
}

func (d *Dialer) auth(addr n.Addr, conn *blend.Conn) error {
	defer conn.Close()
	if err := conn.Write(&HelloMsg{
		SourceAddr: d.dialback,
		TargetAddr: addr,
	}); err != nil {
		return err
	}
	msg, err := conn.Read()
	if err != nil {
		return err
	}
	switch q := msg.(type) {
	case *WelcomeMsg:
		return nil
	case *RejectMsg:
		return errors.NewError("dial rejected by remote (%s)", errors.Unpack(q.Err))
	}
	return errors.NewError("unknown welcome response")
}
