// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package xor implements a nearest-neighbor data structure for the XOR-metric
package xor

import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

// Key represents a point in the XOR-space
type Key uint64

// Key implements interface Item
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

// Item is any type that has an XOR-space Key
type Item interface {
	Key() Key
}

// Metric is an XOR-metric space that supports point addition and nearest neighbor (NN) queries.
// The zero value is an empty metric space.
type Metric struct {
	Item
	sub [2]*Metric
	n   int // Number of items (not nodes) in the subtree of and including this node
}

var ErrDup = errors.New("duplicate point")

// Iterate calls f on each node of the XOR-tree that has a non-nil item.
func (m *Metric) Iterate(f func(Item)) {
	if m.Item != nil {
		f(m.Item)
	}
	if m.sub[0] != nil {
		m.sub[0].Iterate(f)
	}
	if m.sub[1] != nil {
		m.sub[1].Iterate(f)
	}
}

func (m *Metric) Dump() []Item {
	var result []Item
	m.Iterate(func(item Item) {
		result = append(result, item)
	})
	return result
}

// Copy returns a deep copy of the metric
func (m *Metric) Copy() *Metric {
	m_ := &Metric{
		Item: m.Item,
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
	if m.Item != nil {
		m.n++
	}
}

// Add adds the item to the metric. It returns the smallest number of
// significant bits that distinguish this item from the rest in the metric.
func (m *Metric) Add(item Item) (level int, err error) {
	return m.add(item, 0)
}

func (m *Metric) add(item Item, r int) (bottom int, err error) {
	defer m.calcSize()
	if m.Item == nil {
		if m.sub[0] == nil && m.sub[1] == nil {
			// This is an empty leaf node
			m.Item = item
			return r, nil
		}
		// This is an intermediate node
		return m.forward(item, r)
	}
	// This is a non-empty leaf node
	if m.Item.Key() == item.Key() {
		return r, ErrDup
	}
	if _, err = m.forward(m.Item, r); err != nil {
		panic("¢")
	}
	m.Item = nil
	bottom, err = m.forward(item, r)
	if err != nil {
		panic("¢")
	}
	return bottom, err
}

func (m *Metric) forward(item Item, r int) (bottom int, err error) {
	j := item.Key().Bit(r)
	if m.sub[j] == nil {
		m.sub[j] = &Metric{}
	}
	return m.sub[j].add(item, r+1)
}

// Remove removes an item with id from the metric, if present.
// It returns the removed item, or nil if non present.
func (m *Metric) Remove(id Key) Item {
	item, _ := m.remove(id, 0)
	return item
}

func (m *Metric) remove(id Key, r int) (Item, bool) {
	defer m.calcSize()
	if m.Item != nil {
		if m.Item.Key() == id {
			item := m.Item
			m.Item = nil
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
func (m *Metric) Nearest(pivot Key, k int) []Item {
	return m.nearest(pivot, k, 0)
}

func (m *Metric) nearest(pivot Key, k int, r int) []Item {
	if k == 0 {
		return nil
	}
	if m.Item != nil {
		return []Item{m.Item}
	}
	var result []Item
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
