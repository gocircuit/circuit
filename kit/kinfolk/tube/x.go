// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tube

import (
	"time"

	"github.com/gocircuit/circuit/kit/kinfolk"
	"github.com/gocircuit/circuit/use/circuit"
)

// XTube
type XTube struct {
	t *Tube
}

func init() {
	circuit.RegisterValue(XTube{})
}

// XID returns an XID pointing to this Tube table
func (x XTube) XID() kinfolk.FolkXID {
	return x.t.permXID
}

func (x XTube) Subscribe(downXID kinfolk.FolkXID, upsync []*Record) []*Record {
	// log.Printf("xtube subscribing")
	// defer func() {
	// 	log.Printf("xtube subscribed\n%s", x.t.Dump())
	// }()

	if downXID.ID == x.t.permXID.ID {
		panic("x")
	}
	//
	x.t.Lock()
	defer x.t.Unlock()

	x.t.downstream.Open(kinfolk.XID(downXID)) // Add peer to downstream list
	go x.t.BulkWrite(upsync)                  // Catch up to peer after we return (and peer unlocks itself)
	return x.t.bulkRead()
}

func (x XTube) Write(key string, rev Rev, value interface{}) bool {
	// log.Printf("xtube writing (%s,%d,%v)", key, rev, value)
	// defer func() {
	// 	log.Printf("xtube wrote %s\n%s", key, x.t.Dump())
	// }()
	return x.t.Write(key, rev, value)
}

func (x XTube) BulkWrite(bulk []*Record) {
	if len(bulk) == 0 {
		return
	}
	x.t.BulkWrite(bulk)
}

func (x XTube) Scrub(key string, notAfterRev Rev, notAfterUpdated time.Time) {
	x.t.Scrub(key, notAfterRev, notAfterUpdated)
}

// YTube
type YTube struct {
	xid kinfolk.FolkXID
}

// Lookup and Forget intentionally omitted. To be called only by local tube user.

func (y YTube) Subscribe(downXID kinfolk.FolkXID, upsync []*Record) []*Record {
	defer func() {
		if r := recover(); r != nil {
			// log.Printf("ytube subscribe panic\n%#v\n", r)
		}
	}()
	return y.xid.X.Call("Subscribe", downXID, upsync)[0].([]*Record)
}

func (y YTube) Write(key string, rev Rev, value interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			//log.Printf("ytube write panic\n%#v\n", r)
		}
	}()
	return y.xid.X.Call("Write", key, rev, value)[0].(bool)
}

func (y YTube) BulkWrite(bulk []*Record) {
	defer func() {
		if r := recover(); r != nil {
			//log.Printf("ytube bulk write panic\n%#v\n", r)
		}
	}()
	y.xid.X.Call("BulkWrite", bulk)
}

func (y YTube) Scrub(key string, notAfterRev Rev, notAfterUpdated time.Time) {
	defer func() {
		if r := recover(); r != nil {
			// log.Printf("ytube scrub panic\n%#v\n", r)
		}
	}()
	y.xid.X.Call("Scrub", key, notAfterRev, notAfterUpdated)
}
