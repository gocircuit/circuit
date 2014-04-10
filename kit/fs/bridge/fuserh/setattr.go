// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package fuserh

import (
	"log"

	"github.com/gocircuit/circuit/kit/fs/fuse"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

func (r *RH) setattr(req *fuse.SetattrRequest, hdr *fuse.Header, nodeFID rh.FID) interface{} {
	if req.Valid.Uid() { // Reject unsupported flags
		log.Println("set UID not supported")
		return fuse.EPERM
	}
	if req.Valid.Atime() {
		log.Println("set access time not supported")
		return fuse.EPERM
	}

	// Special case: Linux issues a size change to truncate a file before opening it OTRUNC.
	// Synthetic file servers (e.g., plumber) honor open(OTRUNC) but not wstat.
	var err error
	if req.Valid.Size() && req.Size == 0 {
		var fid2 rh.FID
		fid2, err = nodeFID.Walk(nil) // Clone the unopened node FID
		if err != nil {
			return RHError{err}.FUSE()
		}
		var flags = rh.Flag{
			Attr: rh.WriteOnly,
			Deny: true,
			// IsUnix: true,
			// Unix:   os.O_WRONLY,
		}
		flags.Truncate = true
		if err = fid2.Open(flags, nil); err != nil { // Open it for truncation
			return RHError{err}.FUSE()
		}
		fid2.Clunk() // Then clunk it
	}

	var fid2 rh.FID // Ensure an open FID for the wstat target file
	if req.Valid.Handle() {
		if fid2 = r.handle.Lookup(req.Handle); fid2 == nil {
			return fuse.ENOENT
		}
	} else {
		// Currently, wstat needs an open FID by the RH spec. (Should it?)
		// If a valid open FID handle was not provided, we'll just open a new cloned FID,
		// without worrying about clunking it. This works correctly within the circuit,
		// as unused FIDs will be clunked during garbage collection.
		fid2, err = nodeFID.Walk(nil) // Clone the unopened node FID
		if err != nil {
			return RHError{err}.FUSE()
		}
		// To keep it simple, we require write permissions to the file, if it is to be wstat-ed.
		// However, this is an overkill as the current RH spec has weaker or incomparable requirements
		// in some cases. E.g. if only the name is being changed, write permissions to the parent directory,
		// but not the file, are needed.
		// TODO(petar): Considering the many different things that wstat can change — contents, meta, name —
		// perhaps its semantics should be spread out to multiple simpler API calls.
		var flag = rh.Flag{
			Attr: rh.WriteOnly,
			Deny: true,
			// IsUnix: true,
			// Unix:   os.O_WRONLY,
		}
		if err = fid2.Open(flag, nil); err != nil {
			return RHError{err}.FUSE()
		}
	}
	// fid2 is now an open-for-writing FID, ready for use by wstat,
	// and not to be closed at the end of the case block.

	// Set meta
	if err := nodeFID.Wstat(sattrWdir(req)); err != nil {
		return RHError{err}.FUSE()
	}

	dir, err := nodeFID.Stat() // Read the new stats and return them
	if err != nil {
		return RHError{err}.FUSE()
	}
	//
	attr, attrValid, _ := fuseAttr(dir)
	return &fuse.SetattrResponse{
		AttrValid: attrValid, // how long Attr can be cached
		Attr:      attr,      // file attributes
	}
}
