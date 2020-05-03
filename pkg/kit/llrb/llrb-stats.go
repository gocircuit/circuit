// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

// GetHeight() returns an item in the tree with key @key, and it's height in the tree
func (t *Tree) GetHeight(key Item) (result Item, depth int) {
	return t.getHeight(t.root, key)
}

func (t *Tree) getHeight(h *Node, item Item) (Item, int) {
	if h == nil {
		return nil, 0
	}
	if t.less(item, h.Item) {
		result, depth := t.getHeight(h.Left, item)
		return result, depth + 1
	}
	if t.less(h.Item, item) {
		result, depth := t.getHeight(h.Right, item)
		return result, depth + 1
	}
	return h.Item, 0
}

// HeightStats() returns the average and standard deviation of the height
// of elements in the tree
func (t *Tree) HeightStats() (avg, stddev float64) {
	av := &avgVar{}
	heightStats(t.root, 0, av)
	return av.GetAvg(), av.GetStdDev()
}

func heightStats(h *Node, d int, av *avgVar) {
	if h == nil {
		return
	}
	av.Add(float64(d))
	if h.Left != nil {
		heightStats(h.Left, d+1, av)
	}
	if h.Right != nil {
		heightStats(h.Right, d+1, av)
	}
}
