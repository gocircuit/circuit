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

// HandleTable
type HandleTable struct {
	sync.Mutex
	id map[fuse.HandleID]rh.FID // FUSE handle ID -> open FID
	n  fuse.HandleID
}

func (h *HandleTable) Init() {
	h.id = make(map[fuse.HandleID]rh.FID)
}

func (h *HandleTable) Size() int {
	h.Lock()
	defer h.Unlock()
	return len(h.id)
}

func (h *HandleTable) Lookup(id fuse.HandleID) rh.FID {
	h.Lock()
	defer h.Unlock()
	return h.id[id]
}

func (h *HandleTable) Add(fid rh.FID) fuse.HandleID {
	h.Lock()
	defer h.Unlock()
	h.n++ // zero is not a valid id
	h.id[h.n] = fid
	return h.n
}

func (h *HandleTable) Scrub(id fuse.HandleID) rh.FID {
	h.Lock()
	defer h.Unlock()
	fid, ok := h.id[id]
	if !ok {
		return nil
	}
	delete(h.id, id)
	return fid
}
