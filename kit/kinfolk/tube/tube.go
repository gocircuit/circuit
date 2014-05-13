// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tube

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/gocircuit/circuit/kit/kinfolk"
	"github.com/gocircuit/circuit/kit/lang"
	"github.com/gocircuit/circuit/use/circuit"
)

// Tube is a folk data structure that maintains a key-value set sorted by key
type Tube struct {
	permXID kinfolk.FolkXID // XID to this tube
	folk    *kinfolk.Folk   // Folk interface of this tube to the kin system
	sync.Mutex
	lookup     map[string]int
	table      []*Record
	downstream *kinfolk.Rotor // Rotor of YTubes
}

// Record
type Record struct {
	Key     string
	Rev     Rev
	Value   interface{}
	Updated time.Time
}

//
type Rev uint64

func New(kin *kinfolk.Kin, topic string) *Tube {
	t := &Tube{
		lookup:     make(map[string]int),
		downstream: kinfolk.NewRotor(),
	}
	t.permXID = kinfolk.FolkXID{
		X:  circuit.PermRef(XTube{t}),
		ID: lang.ComputeReceiverID(t),
	}
	t.folk = kin.Attach(topic, t.permXID)
	go t.loop()
	return t
}

// Dump returns a textual representation of the contents of this Tube table
func (t *Tube) Dump() string {
	t.Lock()
	defer t.Unlock()
	return t.dump()
}

func (t *Tube) dump() string {
	var w bytes.Buffer
	for _, r := range t.table {
		fmt.Fprintf(&w, "%s––(%d)––>%v\n", r.Key, r.Rev, r.Value)
	}
	return w.String()
}

// loop processes joining peers
func (t *Tube) loop() {
	for {
		t.superscribe(t.folk.Replenish()) // Get new upstream node
	}
}

// XTube
type XTube struct {
	t *Tube
}

// XID returns an XID pointing to this Tube table
func (x XTube) XID() kinfolk.FolkXID {
	return x.t.permXID
}

func (t *Tube) superscribe(peerXID kinfolk.FolkXID) {
	// log.Printf("tube superscribing %s", yup.xid.String())
	// defer func() {
	// 	log.Printf("tube superscribed %s\n%s", yup.xid.String(), t.Dump())
	// }()

	t.Lock()
	defer t.Unlock()
	yup := YTube{kinfolk.FolkXID(t.downstream.Open(kinfolk.XID(peerXID)))}
	go t.BulkWrite(yup.Subscribe(t.permXID, t.bulkRead()))
}

func (t *Tube) bulkRead() []*Record {
	r := make([]*Record, len(t.table)) // Make an external copy since the table changes continuously
	copy(r, t.table)
	return r
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

// BulkRead returns a listing of all elements of the Tube table
func (t *Tube) BulkRead() []*Record {
	t.Lock()
	defer t.Unlock()
	//
	return t.bulkRead()
}

// The first write must have revision bigger than 0. Otherwise it won't take effect.
// Write will block until the diffusion of the write operation reaches its terminal nodes.
// This simple form of backpressure ensures self-inflicted DDoS in the presence of software bugs.
func (t *Tube) Write(key string, rev Rev, value interface{}) (changed bool) {
	// log.Printf("tube writing (%s,%d,%v)", key, rev, value)
	// defer func() {
	// 	log.Printf("tube wrote, changed=%v\n%s", changed, t.Dump())
	// }()

	t.Lock()
	defer t.Unlock()
	//
	changed = t.write(&Record{
		Key:     key,
		Rev:     rev,
		Value:   value,
		Updated: time.Now(),
	})
	//
	if changed {
		go t.writeSync(key, rev, value)
	}
	return
}

func (x XTube) Write(key string, rev Rev, value interface{}) bool {
	// log.Printf("xtube writing (%s,%d,%v)", key, rev, value)
	// defer func() {
	// 	log.Printf("xtube wrote %s\n%s", key, x.t.Dump())
	// }()
	return x.t.Write(key, rev, value)
}

func (t *Tube) writeSync(key string, rev Rev, value interface{}) {
	var wg sync.WaitGroup
	for _, downXID := range t.downstream.Opened() {
		ydown := YTube{kinfolk.FolkXID(downXID)}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ydown.Write(key, rev, value)
		}()
	}
	wg.Wait()
}

func (t *Tube) write(r *Record) (changed bool) {
	//
	i, ok := t.lookup[r.Key]
	if ok && r.Rev <= t.table[i].Rev {
		return false
	}
	r.Updated = time.Now()
	if ok {
		t.table[i] = r
	} else {
		t.table = append(t.table, r)
		t.lookup[r.Key] = len(t.table) - 1
	}
	return true
}

// bulkWrite writes the records in bulk by pointer, without  copying them
func (t *Tube) bulkWrite(bulk []*Record) {
	//
	changed := make([]*Record, 0, len(bulk))
	for _, r := range bulk {
		if t.write(r) {
			changed = append(changed, r)
		}
	}
	//
	go t.bulkWriteSync(changed)
}

func (t *Tube) bulkWriteSync(changed []*Record) {
	// Records exchanged within and across tubes are immutable, so no lock is necessary
	var wg sync.WaitGroup
	for _, downXID := range t.downstream.Opened() {
		ydown := YTube{kinfolk.FolkXID(downXID)}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ydown.BulkWrite(changed)
		}()
	}
	wg.Wait()
}

func (t *Tube) BulkWrite(bulk []*Record) {
	// log.Printf("tube bulk writing")
	// defer func() {
	// 	log.Printf("tube bulk wrote\n%s", t.Dump())
	// }()

	if len(bulk) == 0 {
		return
	}
	// Copy each record before storing them internally
	for i, r := range bulk {
		var y = *r
		bulk[i] = &y
	}
	//
	t.Lock()
	defer t.Unlock()
	//
	t.bulkWrite(bulk)
}

func (x XTube) BulkWrite(bulk []*Record) {
	println("xtube.BulkWrite")
	if len(bulk) == 0 {
		return
	}
	x.t.BulkWrite(bulk)
}

func (t *Tube) Forget(key string, notAfterRev Rev, notAfterUpdated time.Time) bool {
	t.Lock()
	defer t.Unlock()
	//
	return t.forget(key, notAfterRev, notAfterUpdated)
}

func (t *Tube) forget(key string, notAfterRev Rev, notAfterUpdated time.Time) bool {
	i, ok := t.lookup[key]
	if !ok {
		return false
	}
	r := t.table[i]
	if notAfterRev != 0 && r.Rev > notAfterRev {
		return false
	}
	if !notAfterUpdated.IsZero() && r.Updated.After(notAfterUpdated) {
		return false
	}
	delete(t.lookup, key)
	//
	n := len(t.table)
	t.table[i] = t.table[n-1]
	t.table = t.table[:n-1]
	if i < len(t.table) {
		t.lookup[t.table[i].Key] = i
	}
	return true
}

func (t *Tube) Scrub(key string, notAfterRev Rev, notAfterUpdated time.Time) {
	t.Lock()
	defer t.Unlock()
	//
	if t.forget(key, notAfterRev, notAfterUpdated) {
		go t.scrubSync(key, notAfterRev, notAfterUpdated)
	}
}

func (t *Tube) scrubSync(key string, notAfterRev Rev, notAfterUpdated time.Time) {
	var wg sync.WaitGroup
	for _, downXID := range t.downstream.Opened() {
		ydown := YTube{kinfolk.FolkXID(downXID)}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ydown.Scrub(key, notAfterRev, notAfterUpdated)
		}()
	}
	wg.Wait()
}

func (x XTube) Scrub(key string, notAfterRev Rev, notAfterUpdated time.Time) {
	x.t.Scrub(key, notAfterRev, notAfterUpdated)
}

func (t *Tube) Lookup(key string) *Record {
	t.Lock()
	defer t.Unlock()
	//
	i, present := t.lookup[key]
	if !present {
		return nil
	}
	var r = *t.table[i] // Return a copy of the record
	return &r
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

// Init
func init() {
	circuit.RegisterValue(&Tube{}) // In order to be able to compute receiver ID
	circuit.RegisterValue(XTube{})
}
