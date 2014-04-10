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
)

func (r *RH) lookup(req *fuse.LookupRequest, nodeFID rh.FID) interface{} {
	//
	fid2, err := nodeFID.Walk([]string{req.Name})
	if err != nil {
		return RHError{err}.FUSE()
	}
	//
	dir2, err := fid2.Stat()
	if err != nil {
		return RHError{err}.FUSE()
	}
	//
	attr2, attrValid2, entryValid2 := fuseAttr(dir2)
	nid, gen := r.node.Add(fid2)
	return &fuse.LookupResponse{
		Node:       nid,
		Generation: gen,
		EntryValid: entryValid2,
		AttrValid:  attrValid2,
		Attr:       attr2,
	}
}

func (r *RH) forget(hdr *fuse.Header, nodeFID rh.FID) interface{} {
	r.node.Scrub(hdr.Node)
	return nil
}
