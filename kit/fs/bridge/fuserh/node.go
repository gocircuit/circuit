// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package fuserh

import (
	"sync"

	"github.com/gocircuit/circuit/kit/fs/fuse"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// NodeTable
type NodeTable struct {
	sync.Mutex
	id map[fuse.NodeID]rh.FID // Slice index is fuse.NodeID
	n  fuse.NodeID

	// FUSE requires nodeids to have associated generation numbers.
	// If we reuse a nodeid, we have to bump the generation number to
	// guarantee that the nodeid,gen combination is never reused.
	gen uint64 // FUSE generation
}

func (n *NodeTable) Init() {
	n.id = make(map[fuse.NodeID]rh.FID)
	n.gen = 1
}

func (n *NodeTable) Size() int {
	n.Lock()
	defer n.Unlock()
	return len(n.id)
}

func (n *NodeTable) Lookup(nid fuse.NodeID) rh.FID {
	n.Lock()
	defer n.Unlock()
	return n.id[nid]
}

func (n *NodeTable) Add(fid rh.FID) (nid fuse.NodeID, gen uint64) {
	n.Lock()
	defer n.Unlock()
	n.n++ // zero is not a valid nid
	n.id[n.n] = fid
	return n.n, n.gen // We don't increment the generation, because node IDs are never reused.
}

func (n *NodeTable) Scrub(nid fuse.NodeID) rh.FID {
	n.Lock()
	defer n.Unlock()
	fid, ok := n.id[nid]
	if !ok {
		return nil
	}
	delete(n.id, nid)
	return fid
}
