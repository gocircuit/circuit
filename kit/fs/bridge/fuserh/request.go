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
)

// An Intr is a channel that signals that a request has been interrupted.
// Being able to receive from the channel means the request has been
// interrupted.
type Intr chan struct{}

func (Intr) String() string { return "fuse/rh.Intr" }

// RequestTable
type RequestTable struct {
	sync.Mutex
	serving map[fuse.RequestID]*Request // FUSE requests currently being served
}

type Request struct {
	tab     *RequestTable
	Request fuse.Request
	Intr    Intr
	sync.Mutex
	intr chan<- struct{}
}

func (q *Request) Interrupt() {
	if q == nil {
		return
	}
	q.Lock()
	defer q.Unlock()
	if q.intr != nil {
		close(q.intr)
	}
	q.intr = nil
}

// Init
func (r *RequestTable) Init() {
	r.serving = make(map[fuse.RequestID]*Request)
}

func (r *RequestTable) Size() int {
	r.Lock()
	defer r.Unlock()
	return len(r.serving)
}

func (r *RequestTable) Add(req fuse.Request) *Request {
	intr := make(Intr, 1)
	q := &Request{tab: r, Request: req, Intr: intr, intr: (chan<- struct{})(intr)}
	//if !IsJunkRequest(req) {
	//	Debugf("<-%s", req)
	//}
	r.Lock()
	defer r.Unlock()
	hdr := req.Hdr()
	if r.serving[hdr.ID] != nil {
		// This happens with OSXFUSE.  Assume it's okay and
		// that we'll never see an interrupt for this one.
		// Otherwise everything wedges.  TODO: Report to OSXFUSE?
		q.Intr = nil
	} else {
		r.serving[hdr.ID] = q
	}
	return q
}

func (r *RequestTable) Lookup(id fuse.RequestID) *Request {
	r.Lock()
	defer r.Unlock()
	return r.serving[id]
}

func (q *Request) Clunk(resp interface{}) {
	hdr := q.Request.Hdr()
	//if !IsJunkRequest(q.Request) {
	//	Debugf("-> %#x %v", hdr.ID, resp)
	//}
	q.tab.Lock()
	defer q.tab.Unlock()
	delete(q.tab.serving, hdr.ID)
}
