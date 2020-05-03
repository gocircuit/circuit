// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package prof implements internal profiling data structures
package prof

import (
	"sync"
	"time"

	"github.com/hoijui/circuit/pkg/kit/stat"
)

// Profile keeps various load-related statistics for a worker
type Profile struct {
	rlk        sync.Mutex
	replyTotal sketch
	replyProc  map[string]*sketch

	clk       sync.Mutex
	callTotal sketch
	callProc  map[string]*sketch
}

type sketch struct {
	Begin    int64
	End      int64
	Abort    int64
	Dur      stat.Moment
	AbortDur stat.Moment
}

// NewProfile creates a new Profile instance
func New() *Profile {
	return &Profile{
		replyProc: make(map[string]*sketch),
		callProc:  make(map[string]*sketch),
	}
}

func (p *Profile) Stat() *WorkerStat {
	r := &WorkerStat{}

	p.rlk.Lock()
	r.ReplyProc = make(map[string]*Stat)
	for name, sk := range p.replyProc {
		r.ReplyProc[name] = &Stat{
			Type:           "reply",
			Begin:          sk.Begin,
			End:            sk.End,
			Abort:          sk.Abort,
			DurAvg:         sk.Dur.Average(),
			DurStdDev:      sk.Dur.StdDev(),
			AbortDurAvg:    sk.AbortDur.Average(),
			AbortDurStdDev: sk.AbortDur.StdDev(),
		}
	}
	r.ReplyTotal = &Stat{
		Type:           "reply",
		Begin:          p.replyTotal.Begin,
		End:            p.replyTotal.End,
		Abort:          p.replyTotal.Abort,
		DurAvg:         p.replyTotal.Dur.Average(),
		DurStdDev:      p.replyTotal.Dur.StdDev(),
		AbortDurAvg:    p.replyTotal.AbortDur.Average(),
		AbortDurStdDev: p.replyTotal.AbortDur.StdDev(),
	}
	p.rlk.Unlock()

	p.clk.Lock()
	r.CallProc = make(map[string]*Stat)
	for name, sk := range p.callProc {
		r.CallProc[name] = &Stat{
			Type:           "call",
			Begin:          sk.Begin,
			Abort:          sk.Abort,
			End:            sk.End,
			DurAvg:         sk.Dur.Average(),
			DurStdDev:      sk.Dur.StdDev(),
			AbortDurAvg:    sk.AbortDur.Average(),
			AbortDurStdDev: sk.AbortDur.StdDev(),
		}
	}
	r.CallTotal = &Stat{
		Type:           "call",
		Begin:          p.callTotal.Begin,
		End:            p.callTotal.End,
		Abort:          p.callTotal.Abort,
		DurAvg:         p.callTotal.Dur.Average(),
		DurStdDev:      p.callTotal.Dur.StdDev(),
		AbortDurAvg:    p.callTotal.AbortDur.Average(),
		AbortDurStdDev: p.callTotal.AbortDur.StdDev(),
	}
	p.clk.Unlock()

	return r
}

// BeginReply starts a stop watch, measuring the duration of an execution.
// It returns the stop watch interface.
func (p *Profile) BeginReply(proc string) StopWatch {
	p.rlk.Lock()
	defer p.rlk.Unlock()
	// Update total sketch
	p.replyTotal.Begin++
	// Update specific sketch
	p.replyGet(proc).Begin++

	return NewStopWatch(
		func(dur time.Duration) { p.stopReply(proc, dur) },
		func(dur time.Duration) { panic("replies cannot abort") },
	)
}

// BeginCall starts a stop watch, measuring the duration of a call.
// It returns the stop watch interface.
func (p *Profile) BeginCall(proc string) StopWatch {
	p.clk.Lock()
	defer p.clk.Unlock()
	// Update total sketch
	p.callTotal.Begin++
	// Update specific sketch
	p.callGet(proc).Begin++

	return NewStopWatch(
		func(dur time.Duration) { p.stopCall(proc, dur) },
		func(dur time.Duration) { p.abortCall(proc, dur) },
	)
}

// replyGet must be called under a lock
func (p *Profile) replyGet(proc string) *sketch {
	sk, present := p.replyProc[proc]
	if !present {
		sk = &sketch{}
		p.replyProc[proc] = sk
	}
	return sk
}

// callGet must be called under a lock
func (p *Profile) callGet(proc string) *sketch {
	sk, present := p.callProc[proc]
	if !present {
		sk = &sketch{}
		p.callProc[proc] = sk
	}
	return sk
}
