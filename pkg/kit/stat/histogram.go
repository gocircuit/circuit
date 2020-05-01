// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package stat

// Histogram is a simple static histogram structure.
// It consumes a stream of floating point numbers and can output a histogram at any time.
type Histogram struct {
	min, max float64
	width    float64
	bin      []*Bin
}

type Bin struct {
	X      float64 // Starting value. Bin is [X, X+width)
	Weight float64
}

// NewHistogram creates and returns a new histogram object.
func NewHistogram(min, max float64, n int) *Histogram {
	h := &Histogram{
		min:   min,
		max:   max,
		width: (max - min) / float64(n),
		bin:   make([]*Bin, n),
	}
	if h.width < 0 {
		panic("negative histogram bin width")
	}
	for i, _ := range h.bin {
		h.bin[i] = &Bin{X: min + float64(i)*h.width}
	}
	return h
}

// Width returns the histogram bin width in the value domain
func (h *Histogram) Width() float64 {
	return h.width
}

// Put places the sample with value x and given weight (must be positive) into the histogram.
func (h *Histogram) Put(x float64, weight float64) {
	if x <= h.min {
		h.bin[0].Weight += weight
		return
	}
	if x >= h.max {
		h.bin[len(h.bin)-1].Weight += weight
	}
	h.bin[min(len(h.bin)-1, int((x-h.min)/h.width))].Weight += weight
}

// Histogram returns the current histogram.
func (h *Histogram) Histogram() []*Bin {
	return h.bin
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
