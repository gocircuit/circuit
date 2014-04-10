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

func (r *RH) release(req *fuse.ReleaseRequest, hdr *fuse.Header) interface{} {
	fid2 := r.handle.Scrub(req.Handle)
	if fid2 == nil {
		Debugf("-> %#x %v", hdr.ID, fuse.ESTALE)
		return fuse.ESTALE
	}
	if err := fid2.Clunk(); err != nil {
		return RHError{err}.FUSE()
	}
	return nil
}

func (r *RH) rename(req *fuse.RenameRequest, nodeFID rh.FID) interface{} {
	fid2 := r.node.Lookup(req.NewDir)
	if fid2 == nil {
		return fuse.EPERM
	}
	if err := nodeFID.Move(fid2, req.NewName); err != nil {
		return RHError{err}.FUSE()
	}
	return nil
}
