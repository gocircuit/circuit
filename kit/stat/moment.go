// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package stat implements facilities for storing distribution sketches and updating them
package stat

import (
	"math"
)

// Moment is a streaming sketch data structure.
// It keeps track of moment statistics for an incoming stream of floating point numbers.
type Moment struct {
	sum    float64
	sumAbs float64
	sumSq  float64
	min    float64
	max    float64
	weight float64
	n      int64
}

// Init initializes the sketch.
func (x *Moment) Init() {
	x.sum, x.sumAbs, x.sumSq, x.min, x.max, x.weight, x.n = 0, 0, 0, math.NaN(), math.NaN(), 0, 0
}

// Add adds sample with weight one.
func (x *Moment) Add(sample float64) {
	x.AddWeighted(sample, 1)
}

// AddWeighted adds sample with the given weight.
func (x *Moment) AddWeighted(sample float64, weight float64) {
	x.n++
	x.sum += sample * weight
	x.sumAbs += math.Abs(sample * weight)
	x.sumSq += sample * sample * weight
	x.weight += weight
	if math.IsNaN(x.min) || sample < x.min {
		x.min = sample
	}
	if math.IsNaN(x.max) || sample > x.max {
		x.max = sample
	}
}

// IsEmpty returns true if no samples have been added to this sketch yet.
func (x *Moment) IsEmpty() bool {
	return x.n == 0
}

// Count returns the number of samples added to this sketch.
func (x *Moment) Count() int64 {
	return x.n
}

// Weight returns the total sample weight added to this sketch.
func (x *Moment) Weight() float64 {
	return x.weight
}

// Mass returns the sum of absolute values of the sample-weight product of all added samples.
func (x *Moment) Mass() float64 {
	return x.sumAbs
}

// Average returns the weighted average of samples added to the sketch.
func (x *Moment) Average() float64 {
	return x.Moment(1)
}

// Variance returns the weighted variance of the samples added to this sketch.
func (x *Moment) Variance() float64 {
	m1 := x.Moment(1)
	return x.Moment(2) - m1*m1
}

// StdDev returns the weighted standard deviation of the samples added to this sketch.
func (x *Moment) StdDev() float64 {
	return math.Sqrt(x.Variance())
}

// Min returns the smallest sample added to this sketch.
func (x *Moment) Min() float64 {
	return x.min
}

// Ma returns the largest sample added to this sketch.
func (x *Moment) Max() float64 {
	return x.max
}

// Moment returns the weighted k-th moment of the samples added to this sketch.
func (x *Moment) Moment(k float64) float64 {
	switch k {
	case 0:
		return 1
	case 1:
		return x.sumAbs / x.weight
	case 2:
		return x.sumSq / x.weight
	}
	if math.IsInf(k, 1) {
		return x.max
	}
	panic("not yet supported")
}
