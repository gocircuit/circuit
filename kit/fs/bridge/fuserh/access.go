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

func (r *RH) access(req *fuse.AccessRequest, nodeFID rh.FID) interface{} {
	return fuse.ENOSYS
	//
	// Compute RH flag corresponding to access mask (see man page for access)
	// Clone the FID
	/*
		fid2, err := nodeFID.Walk(nil)
		if err != nil {
			return RHError{err}.FUSE()
		}
		if err = fid2.Open(accessMask(req.Mask).RH()); err != nil {
			return RHError{err}.FUSE()
		}
		fid2.Clunk()
		return nil
	*/
}
