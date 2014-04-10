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

func (r *RH) getattr(nodeFID rh.FID) interface{} {
	dir, err := nodeFID.Stat()
	if err != nil {
		return RHError{err}.FUSE()
	}
	//
	attr, attrValid, _ := fuseAttr(dir)
	return &fuse.GetattrResponse{
		AttrValid: attrValid, // how long Attr can be cached
		Attr:      attr,      // file attributes
	}
}
