// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package prof

import (
	"time"
)

// StopWatch represents a stop watch
type StopWatch interface {
	Stop()
	Abort()
}

// stopWatch is a stop watch for function execution duration
type stopWatch struct {
	t0    time.Time
	stop  func(time.Duration)
	abort func(time.Duration)
}

func NewStopWatch(stopFunc, abortFunc func(time.Duration)) StopWatch {
	return &stopWatch{t0: time.Now(), stop: stopFunc, abort: abortFunc}
}

func (x *stopWatch) Stop() {
	dur := time.Now().Sub(x.t0)
	x.stop(dur)
}

func (x *stopWatch) Abort() {
	dur := time.Now().Sub(x.t0)
	x.abort(dur)
}

// Profile hooks for stopwatch

func (p *Profile) stopReply(key string, dur time.Duration) {
	p.rlk.Lock()
	defer p.rlk.Unlock()
	// Add totals
	p.replyTotal.End++
	p.replyTotal.Dur.Add(float64(dur))
	// Add specifics
	sk := p.replyGet(key)
	sk.End++
	sk.Dur.Add(float64(dur))
}

func (p *Profile) stopCall(key string, dur time.Duration) {
	p.clk.Lock()
	defer p.clk.Unlock()
	// Add totals
	p.callTotal.End++
	p.callTotal.Dur.Add(float64(dur))
	// Add specifics
	sk := p.callGet(key)
	sk.End++
	sk.Dur.Add(float64(dur))
}

func (p *Profile) abortCall(key string, dur time.Duration) {
	p.clk.Lock()
	defer p.clk.Unlock()
	// Add totals
	p.callTotal.Abort++
	p.callTotal.AbortDur.Add(float64(dur))
	// Add specifics
	sk := p.callGet(key)
	sk.Abort++
	sk.AbortDur.Add(float64(dur))
}
