// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package fuserh

import (
	"github.com/gocircuit/circuit/kit/fs/fuse"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/bridge/rhunix"
)

func (r *RH) open(q *Request, req *fuse.OpenRequest, nodeFID rh.FID) interface{} {
	//Debugf("FUSE—>LOOPFS %v", req)
	// Clone the FID
	fid2, err := nodeFID.Walk(nil)
	if err != nil {
		return RHError{err}.FUSE()
	}
	// Check against req.Dir
	dir2, err := fid2.Stat()
	if err != nil {
		return RHError{err}.FUSE()
	}
	switch {
	case req.Dir && !dir2.IsDir(), !req.Dir && dir2.IsDir():
		return fuse.EPERM
	}
	// Capture any interrupts targeted to this request and fire a prompt to the RH FID
	eofn, intr := make(chan struct{}, 1), make(chan rh.Prompt, 1)
	defer close(eofn)
	go func() {
		select {
		case <-eofn:
		case <-q.Intr:
			intr <- nil
		}
		close(intr)
	}()
	//
	// Ignore req.Mode. Does not apply in RH or 9P.
	rhflag := rhunix.UnixFlag(fuseOpenFlag(req.Flags).Unix()).RH()
	if err = fid2.Open(rhflag, intr); err != nil {
		return RHError{err}.FUSE()
	}
	var flags fuse.OpenFlags
	if dir2.Mode.Attr == rh.ModeIO {
		flags = fuse.OpenDirectIO
	}
	h := r.handle.Add(fid2)
	return &fuse.OpenResponse{
		Handle: h,
		Flags:  flags,
	}
}

func (r *RH) create(req *fuse.CreateRequest, nodeFID rh.FID) interface{} {
	mp := rhunix.UnixMode(req.Mode).RH()
	rhflag := rhunix.UnixFlag(fuseOpenFlag(req.Flags).Unix()).RH()
	fid2, err := nodeFID.Create(req.Name, rhflag, mp.Mode, mp.Perm)
	if err != nil {
		return RHError{err}.FUSE()
	}
	dir2, err := fid2.Stat()
	if err != nil {
		return RHError{err}.FUSE()
	}
	attr, attrValid, entryValid := fuseAttr(dir2)
	nid, gen := r.node.Add(fid2)
	var flags fuse.OpenFlags
	if dir2.Mode.Attr == rh.ModeIO {
		flags = fuse.OpenDirectIO
	}
	return &fuse.CreateResponse{
		LookupResponse: fuse.LookupResponse{
			Node:       nid,
			Generation: gen,
			EntryValid: entryValid,
			AttrValid:  attrValid,
			Attr:       attr,
		},
		OpenResponse: fuse.OpenResponse{
			Handle: r.handle.Add(fid2),
			Flags:  flags,
		},
	}
}
