// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package stat

import (
	"time"
)

// TimeSampler is a facility for collection stopwatch statistics over multiple experiments.
// TimeSampler is not synchronized. Only one measurement can take place at a time.
type TimeSampler struct {
	m  Moment
	t0 *time.Time
}

// Init initializes the time sampler.
func (x *TimeSampler) Init() {
	x.m.Init()
	x.t0 = nil
}

// Start initiates a new measurement.
func (x *TimeSampler) Start() {
	if x.t0 != nil {
		panic("previous sample not completed")
	}
	t0 := time.Now()
	x.t0 = &t0
}

// Stop ends an experiment and records the elapsed time as a sample in an underlying moment sketch.
func (x *TimeSampler) Stop() {
	t1 := time.Now()
	diff := t1.Sub(*x.t0)
	x.t0 = nil
	x.m.Add(float64(diff))
}

// Moment returns the underlying moment sketch.
func (x *TimeSampler) Moment() *Moment {
	return &x.m
}

// Average returns the average experiment time.
func (x *TimeSampler) Average() float64 {
	return x.m.Average()
}

// StdDev returns the standard deviation across all experiments.
func (x *TimeSampler) StdDev() float64 {
	return x.m.StdDev()
}
