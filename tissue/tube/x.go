// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tube

import (
	"time"

	"github.com/gocircuit/circuit/tissue"
	"github.com/gocircuit/circuit/use/circuit"
)

// XTube is the interface to a tube, given to the tube's upstream (cross-circiuit) peering tubes.
type XTube struct {
	t *Tube
}

func init() {
	circuit.RegisterValue(XTube{})
}

// Avatar returns an Avatar pointing to this tube
func (x XTube) Avatar() tissue.FolkAvatar {
	return x.t.av
}

// Write…
func (x XTube) Write(key string, rev Rev, value interface{}) bool {
	// log.Printf("xtube writing (%s,%d,%v)", key, rev, value)
	// defer func() {
	// 	log.Printf("xtube wrote %s\n%s", key, x.t.Dump())
	// }()
	return x.t.Write(key, rev, value)
}

func (x XTube) BulkWrite(bulk []*Record) {
	x.t.BulkWrite(bulk)
}

func (x XTube) Scrub(key string, notAfterRev Rev, notAfterUpdated time.Time) {
	x.t.Scrub(key, notAfterRev, notAfterUpdated)
}

// YTube
type YTube struct {
	av tissue.FolkAvatar
}

// Lookup and Forget intentionally omitted. To be called only by local tube user.

func (y YTube) Write(key string, rev Rev, value interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			//log.Printf("ytube write panic\n%#v\n", r)
		}
	}()
	return y.av.X.Call("Write", key, rev, value)[0].(bool)
}

func (y YTube) BulkWrite(bulk []*Record) {
	defer func() {
		if r := recover(); r != nil {
			//log.Printf("ytube bulk write panic\n%#v\n", r)
		}
	}()
	y.av.X.Call("BulkWrite", bulk)
}

func (y YTube) Scrub(key string, notAfterRev Rev, notAfterUpdated time.Time) {
	defer func() {
		if r := recover(); r != nil {
			// log.Printf("ytube scrub panic\n%#v\n", r)
		}
	}()
	y.av.X.Call("Scrub", key, notAfterRev, notAfterUpdated)
}
