// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package xor implements a nearest-neighbor data structure for the XOR-metric.
package xor

import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

// Key represents a point in the 64-bit XOR-space.
type Key uint64

// Key implements interface Point
func (id Key) Key() Key {
	return id
}

// Bit returns the k-th bit from metric point of view. The zero-th bit is most significant.
func (id Key) Bit(k int) int {
	return int((id >> uint(k)) & 1)
}

// String returns a textual representation of the id
func (id Key) String() string {
	return fmt.Sprintf("%016x", id)
}

// String returns a textual representation of the id, truncated to the k MSBs.
func (id Key) ShortString(k uint) string {
	shift := uint(8*unsafe.Sizeof(id)) - k
	return fmt.Sprintf("%0"+strconv.Itoa(int(k))+"b", ((id << shift) >> shift))
}

// Point is any type that has an XOR-space Key
type Point interface {
	Key() Key
}

// Proximity equals the number of contiguous bits, starting from 
// the metrically most significant one, that x and y share.
func Proximity(x, y Point) int {
	const nbits = 64
	xy := x.Key() ^ y.Key()
	for i := 0; i < nbits; i++ {
		if xy.Bit(i) != 0 {
			return i
		}
	}
	return nbits
}

// Metric is an XOR-metric space that supports point addition and nearest neighbor (NN) queries.
// The zero value is an empty metric space.
type Metric struct {
	Point
	sub [2]*Metric
	n   int // Number of items (not nodes) in the subtree of and including this node
}

var ErrDup = errors.New("duplicate point")

// Iterate calls f on each node of the XOR-tree that has a non-nil item.
func (m *Metric) Iterate(f func(Point)) {
	if m.Point != nil {
		f(m.Point)
	}
	if m.sub[0] != nil {
		m.sub[0].Iterate(f)
	}
	if m.sub[1] != nil {
		m.sub[1].Iterate(f)
	}
}

func (m *Metric) Dump() []Point {
	var result []Point
	m.Iterate(func(item Point) {
		result = append(result, item)
	})
	return result
}

// Copy returns a deep copy of the metric
func (m *Metric) Copy() *Metric {
	m_ := &Metric{
		Point: m.Point,
		n:    m.n,
	}
	if m.sub[0] != nil {
		m_.sub[0] = m.sub[0].Copy()
	}
	if m.sub[1] != nil {
		m_.sub[1] = m.sub[1].Copy()
	}
	return m_
}

// Clear removes all points from the metric
func (m *Metric) Clear() {
	*m = Metric{}
}

// Size returns the number of points in the metric
func (m *Metric) Size() int {
	return m.n
}

func (m *Metric) calcSize() {
	m.n = 0
	if m.sub[0] != nil {
		m.n += m.sub[0].n
	}
	if m.sub[1] != nil {
		m.n += m.sub[1].n
	}
	if m.Point != nil {
		m.n++
	}
}

// Add adds the item to the metric. It returns the smallest number of
// significant bits that distinguish this item from the rest in the metric.
func (m *Metric) Add(item Point) (level int, err error) {
	return m.add(item, 0)
}

func (m *Metric) add(item Point, r int) (bottom int, err error) {
	defer m.calcSize()
	if m.Point == nil {
		if m.sub[0] == nil && m.sub[1] == nil {
			// This is an empty leaf node
			m.Point = item
			return r, nil
		}
		// This is an intermediate node
		return m.forward(item, r)
	}
	// This is a non-empty leaf node
	if m.Point.Key() == item.Key() {
		return r, ErrDup
	}
	if _, err = m.forward(m.Point, r); err != nil {
		panic("¢")
	}
	m.Point = nil
	bottom, err = m.forward(item, r)
	if err != nil {
		panic("¢")
	}
	return bottom, err
}

func (m *Metric) forward(item Point, r int) (bottom int, err error) {
	j := item.Key().Bit(r)
	if m.sub[j] == nil {
		m.sub[j] = &Metric{}
	}
	return m.sub[j].add(item, r+1)
}

// Remove removes an item with id from the metric, if present.
// It returns the removed item, or nil if non present.
func (m *Metric) Remove(id Key) Point {
	item, _ := m.remove(id, 0)
	return item
}

func (m *Metric) remove(id Key, r int) (Point, bool) {
	defer m.calcSize()
	if m.Point != nil {
		if m.Point.Key() == id {
			item := m.Point
			m.Point = nil
			return item, true
		}
		return nil, false
	}
	b := id.Bit(r)
	sub := m.sub[b]
	if sub == nil {
		return nil, false
	}
	item, emptied := sub.remove(id, r+1)
	if emptied {
		m.sub[b] = nil
		if m.sub[1-b] == nil {
			return item, true
		}
	}
	return item, false
}

// Nearest returns the k points in the metric that are closest to the pivot.
func (m *Metric) Nearest(pivot Key, k int) []Point {
	return m.nearest(pivot, k, 0)
}

func (m *Metric) nearest(pivot Key, k int, r int) []Point {
	if k == 0 {
		return nil
	}
	if m.Point != nil {
		return []Point{m.Point}
	}
	var result []Point
	b := pivot.Bit(r)
	sub := m.sub[b]
	if sub != nil {
		result = sub.nearest(pivot, k, r+1)
	}
	k -= len(result)
	sub = m.sub[1-b]
	if sub != nil {
		result = append(result, sub.nearest(pivot, k, r+1)...)
	}
	return result
}
