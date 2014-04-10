// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package fuserh provides a FUSE server for an RH file system.
//
// In our modular notation for file system chains,
// this package provides the module FUSE<RH>: A FUSE wrapper for RH.
//
package fuserh

// Much of the RH-to-FUSE conversion logic is adapted from the 9P-to-FUSE logic in:
//	http://swtch.com/usr/local/plan9/src/cmd/9pfuse/main.c

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"runtime/debug"
	"syscall"

	"github.com/gocircuit/circuit/kit/fs/fuse"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/sched/limiter"
)

// RH is a FUSE file system server for a RH file system
type RH struct {
	qmax    int // Maximum number of concurrent requests
	ssn     rh.Session
	conn    *fuse.Conn
	eof     EOF
	fuseEOF EOF
	//
	node    NodeTable
	serving RequestTable
	handle  HandleTable
}

// Mount mounts the file system of the resource hierarchy ssn onto the local directory dir.
func Mount(dir string, ssn rh.Session, qmax int) (rh *RH, err error) {

	// Try connecting to FUSE twice. The first attempt, if failed, has the effect of
	// unmounting orphan mounted directories.
	var conn *fuse.Conn
	for i := 0; i < 2; i++ {
		if conn, err = fuse.Mount(dir); err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	//
	rh = &RH{
		qmax: qmax,
		ssn:  ssn,
		conn: conn,
	}
	rh.eof.Init()
	rh.fuseEOF.Init()
	rh.node.Init()
	rh.serving.Init()
	rh.handle.Init()

	go func() { // Try to unmount the local mount directory when this process dies.
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Kill, os.Interrupt, syscall.SIGSTOP, syscall.SIGPIPE, syscall.SIGHUP, syscall.SIGFPE, syscall.SIGSEGV)
		<-ch
		fuse.Umount(dir)
	}()

	go rh.loop() // Start accepting FUSE requests

	return rh, nil
}

func (r *RH) dump() {
	nnode, _, nhandle := r.node.Size(), r.serving.Size(), r.handle.Size()
	log.Printf("node=%d, handle=%d", nnode, nhandle)
}

// EOF blocks until the FUSE connection is alive.
// It returns the reason for closure.
func (r *RH) EOF() error {
	return r.eof.EOF()
}

// loop listens for FUSE requests until the connection is closed.
func (r *RH) loop() {
	lmtr := limiter.New(r.qmax)
	defer lmtr.Wait()
	rr := newReadRequestChan(r.conn)
	//
	rootFID, err := r.ssn.Walk(nil)
	if err != nil {
		panic("u")
	}
	if rootNID, _ := r.node.Add(rootFID); rootNID != 1 { // Insert root into node table with node ID 1
		panic("u")
	}
	for {
		select {
		case q := <-rr:
			if err, ok := q.(error); ok && err != nil {
				r.eof.Close(err)
				return
			}
			lmtr.Go(func() {
				r.serve(q.(fuse.Request))
			})
		case err := <-r.fuseEOF.Chan():
			r.eof.Close(err)
			return
		}
	}
}

func newReadRequestChan(conn *fuse.Conn) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			q, err := conn.ReadRequest()
			if err != nil {
				ch <- err
				return
			}
			ch <- q
		}
	}()
	return ch
}

// serve responds to an individual FUSE request.
func (r *RH) serve(req fuse.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("fuse/rh serve panic: %v", r)
			debug.PrintStack()
		}
	}()
	// Debugf("◘ %v", req)
	//defer r.dump()
	//Debugf("<= %#T", req)

	// Create request interrupt object
	q := r.serving.Add(req)

	// Call this before responding.
	// After responding is too late: we might get another request
	// with the same ID and be very confused.
	done := q.Clunk

	// Locate FID object, corresponding to Header.Node
	hdr := req.Hdr()
	var nodeFID rh.FID
	if id := hdr.Node; id != 0 {
		nodeFID = r.node.Lookup(hdr.Node) // Locate node,
		if nodeFID == nil {               // Otherwise,
			//Debugln("missing node", id)
			done(fuse.ESTALE)
			req.RespondError(fuse.ESTALE)
			return
		}
	}

	// if !IsJunkRequest(req) {
	// 	Debugf("◘ %#x %v", hdr.ID, req)
	// }

	var x interface{}
	switch req := req.(type) {
	default:
		// Note: To FUSE, ENOSYS means "this server never implements this request."
		// It would be inappropriate to return ENOSYS for other operations in this
		// switch that might only be unavailable in some contexts, not all.
		x = fuse.ENOSYS

	case *fuse.InterruptRequest:
		//Debugf("◘ %#x %v", hdr.ID, req)
		r.serving.Lookup(req.IntrID).Interrupt()
		x = nil

	// Requests to the file system object

	case *fuse.InitRequest:
		x = r.init()

	case *fuse.StatfsRequest:
		x = r.statfs()

	case *fuse.DestroyRequest:
		x = r.destroy()

	case *fuse.MkdirRequest:
		x = r.mkdir(req, nodeFID)

	case *fuse.SymlinkRequest, *fuse.ReadlinkRequest, *fuse.LinkRequest: // Resource hierarchies do not support linking
		x = fuse.ENOSYS

	// Requests acting on nodes

	case *fuse.RemoveRequest:
		x = r.remove(req, nodeFID)

	case *fuse.GetattrRequest:
		x = r.getattr(nodeFID)

	case *fuse.SetattrRequest:
		x = r.setattr(req, hdr, nodeFID)

	case *fuse.LookupRequest:
		x = r.lookup(req, nodeFID)

	case *fuse.AccessRequest:
		x = r.access(req, nodeFID)

	case *fuse.OpenRequest:
		x = r.open(q, req, nodeFID)

	case *fuse.CreateRequest:
		x = r.create(req, nodeFID)

	case *fuse.MknodRequest: // Not supported
		x = fuse.ENOSYS

	case *fuse.ForgetRequest: // Forget a FUSE node ID, returned to a Lookup, Create or Mkdir request
		x = r.forget(hdr, nodeFID)

	case *fuse.GetxattrRequest, *fuse.SetxattrRequest, *fuse.ListxattrRequest, *fuse.RemovexattrRequest: // Not supported, OSX extended file attributes
		x = fuse.ENOSYS

	// Requests acting on open file handles

	case *fuse.ReadRequest:
		x = r.read(q, req, hdr)

	case *fuse.WriteRequest:
		x = r.write(q, req, hdr)

	case *fuse.FlushRequest:
		// (9fuse) Flush is supposed to flush any buffered writes.  Don't use this.
		//
		// Flush is a total crock.  It gets called on close() of a file descriptor
		// associated with this open file.  Some open files have multiple file
		// descriptors and thus multiple closes of those file descriptors.
		// In those cases, Flush is called multiple times.  Some open files
		// have file descriptors that are closed on process exit instead of
		// closed explicitly.  For those files, Flush is never called.
		// Even more amusing, Flush gets called before close() of read-only
		// file descriptors too!
		//
		// This is just a bad idea.
		x = nil

	case *fuse.FsyncRequest: // Not supported
		// (9fuse) Fsync commits file info to stable storage.
		x = fuse.ENOSYS

	case *fuse.ReleaseRequest:
		x = r.release(req, hdr)

	case *fuse.RenameRequest:
		x = r.rename(req, nodeFID)
	}
	//
	done(x)
	// if !IsJunkRequest(req) {
	// 	Debugf("‣ %#x %v", hdr.ID, x)
	// }
	switch z := x.(type) {
	case nil:
		req.(interface {
			Respond()
		}).Respond()
	case fuse.Error:
		//Debugf("• %#x %v => %v", hdr.ID, req, z)
		req.RespondError(z)
	default:
		reflect.ValueOf(req).MethodByName("Respond").Call([]reflect.Value{reflect.ValueOf(z)})
	}
}
