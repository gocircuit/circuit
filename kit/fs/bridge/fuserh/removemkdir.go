// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package fuserh

import (
	"path"

	"github.com/gocircuit/circuit/kit/fs/fuse"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/bridge/rhunix"
)

func (r *RH) mkdir(req *fuse.MkdirRequest, nodeFID rh.FID) interface{} {
	// Get parent dir
	dirpath, file := path.Split(req.Name)
	if dirpath != "" {
		panic("x")
	}
	// var wname = split(dirpath)
	// if len(wname) == 1 && wname[0] == "." {
	// 	wname = nil
	// }
	// fid1, err := r.ssn.Walk(wname)
	// if err != nil {
	// 	return RHError{err}.FUSE()
	// }
	fid1 := nodeFID
	// Create within it
	mp := rhunix.UnixMode(req.Mode).RH()
	mp.Mode.Attr = rh.ModeDir
	var flags = rh.Flag{
		Attr:   rh.ReadWrite,
		Create: true,
	}
	fid2, err := fid1.Create(file, flags, mp.Mode, mp.Perm)
	if err != nil {
		return RHError{err}.FUSE()
	}
	fid3, err := fid2.Walk(nil) // fid2 is open, and
	if err != nil {
		return RHError{err}.FUSE()
	}
	fid2.Clunk() // we need to clunk it, but we need an unopened FID pointing to the new dir. That's fid3.
	dir, err := fid3.Stat()
	if err != nil {
		return RHError{err}.FUSE()
	}
	attr, attrValid, entryValid := fuseAttr(dir)
	nid, gen := r.node.Add(fid3)
	return &fuse.MkdirResponse{
		LookupResponse: fuse.LookupResponse{
			Node:       nid,
			Generation: gen,
			EntryValid: entryValid,
			AttrValid:  attrValid,
			Attr:       attr,
		},
	}
}

func (r *RH) remove(req *fuse.RemoveRequest, nodeFID rh.FID) interface{} {
	child, err := nodeFID.Walk([]string{req.Name})
	if err != nil {
		return RHError{err}.FUSE()
	}
	dir, err := child.Stat()
	if err != nil {
		return RHError{err}.FUSE()
	}
	switch {
	case req.Dir && !dir.IsDir(), !req.Dir && dir.IsDir():
		return fuse.EPERM
	}
	if err = child.Remove(); err != nil {
		return RHError{err}.FUSE()
	}
	return nil
}
