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
	"log"
)

func (r *RH) read(q *Request, req *fuse.ReadRequest, hdr *fuse.Header) (reply interface{}) {
	// Debugf("◘ %#x <- %v", hdr.ID, req)
	// defer func() {
	// 	Debugf("• %#x => %v", hdr.ID, reply)
	// }()
	//
	fid2 := r.handle.Lookup(req.Handle)
	if fid2 == nil {
		Debugf("-> %#x %v", hdr.ID, fuse.ESTALE)
		return fuse.ESTALE
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
	var resp = &fuse.ReadResponse{}
	var err error
	if req.Dir { // Read directory
		chunk, err := fid2.Read(0, 0, rh.Intr(intr))
		if err != nil {
			return RHError{err}.FUSE()
		}
		if chunk == nil {
			resp.Data = nil
			return resp
		}
		rhdir, ok := chunk.(rh.DirChunk)
		if !ok {
			panic("x")
		}
		for _, dir := range rhdir {
			//Debugf("▒ %v", dir)
			de := fuse.Dirent{
				Inode: (*rhunix.RHDir)(dir).Inode(),
				Type:  fuse.DT_Unknown,
				Name:  dir.Name,
			}
			resp.Data = fuse.AppendDirent(resp.Data, de)
		}
		off := int(req.Offset)
		if off >= len(resp.Data) {
			resp.Data = nil
			return resp
		}
		resp.Data = resp.Data[off : off+min(len(resp.Data)-off, req.Size)]
		return resp
	} else { // Read file
		var chunk rh.Chunk
		if chunk, err = fid2.Read(req.Offset, req.Size, rh.Intr(intr)); err != nil {
			return RHError{err}.FUSE()
		}
		if chunk != nil {
			if data, ok := chunk.(rh.ByteChunk); ok {
				resp.Data = data
			} else {
				// A convenient way of returning from a Read.
				log.Printf("Warning: Read did not return ByteChunk - attempting conversion")
				resp.Data = rh.ByteChunk(chunk.([]byte))
			}
		}
		return resp
	}
}

func (r *RH) write(q *Request, req *fuse.WriteRequest, hdr *fuse.Header) (reply interface{}) {
	// Debugf("◘ %#x <- %v", hdr.ID, req)
	// defer func() {
	// 	Debugf("• %#x => %v", hdr.ID, reply)
	// }()
	//
	fid2 := r.handle.Lookup(req.Handle)
	if fid2 == nil {
		return fuse.EBADF
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
	var resp = &fuse.WriteResponse{}
	var err error
	var n int
	if n, err = fid2.Write(req.Offset, rh.ByteChunk(req.Data), rh.Intr(intr)); err != nil {
		return RHError{err}.FUSE()
	}
	resp.Size = n
	return resp
}
