// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tube

import (
	// "log"
	"sync"
	"time"

	"github.com/hoijui/circuit/pkg/kit/lang"
	"github.com/hoijui/circuit/pkg/kit/pubsub"
	"github.com/hoijui/circuit/pkg/tissue"
	"github.com/hoijui/circuit/pkg/use/circuit"
)

// Tube is a folk data structure that maintains a key-value set sorted by key.
// …
type Tube struct {
	av   tissue.FolkAvatar // Avatar to this tube
	folk *tissue.Folk      // Folk interface of this tube to the kin system
	sync.Mutex
	view *View
}

func init() {
	circuit.RegisterValue(&Tube{}) // In order to be able to compute receiver ID
}

// NewTube…
func NewTube(kin *tissue.Kin, topic string) *Tube {
	t := &Tube{view: NewView()}
	t.av = tissue.FolkAvatar{
		X:  circuit.PermRef(XTube{t}),
		ID: lang.ComputeReceiverID(t),
	}
	t.folk = kin.Attach(topic, t.av)
	go func() {
		for {
			// Consume identities of new downstream nodes
			t.superscribe(t.folk.Replenish())
		}
	}()
	return t
}

// NewArrivals returns a subscription for the stream of arriving peer identities.
func (t *Tube) NewArrivals() *pubsub.Subscription {
	return t.view.NewArrivals()
}

// NewDepartures returns a subscription for the stream of departing peer identities.
func (t *Tube) NewDepartures() *pubsub.Subscription {
	return t.view.NewDepartures()
}

func (t *Tube) superscribe(peer tissue.FolkAvatar) {
	// log.Printf("tube superscribing %s", peer.ID.String())
	// defer func() {
	// 	log.Printf("tube superscribed %s\n%s", peer.ID.String(), t.Dump())
	// }()
	t.Lock()
	peek := t.view.Peek()
	t.Unlock()
	(YTube{peer}).BulkWrite(peek) // Broadcast our knowledge to joining downstream node.
}

// BulkRead returns a listing of all elements of the Tube table
func (t *Tube) BulkRead() []*Record {
	t.Lock()
	defer t.Unlock()
	return t.view.Peek()
}

// Write updates the state of our local view for the given key.
// Write will block until the diffusion of the write operation reaches its terminal nodes.
//
// This simple form of backpressure ensures self-inflicted DDoS in the presence of software bugs.
func (t *Tube) Write(key string, rev Rev, value interface{}) (changed bool) {
	// log.Printf("tube writing (%s,%d,%v)", key, rev, value)
	// defer func() {
	// 	log.Printf("tube written to, changed=%v\n%s", changed, t.Dump())
	// }()

	t.Lock()
	defer t.Unlock()
	changed = t.view.Update(&Record{
		Key:     key,
		Rev:     rev,
		Value:   value,
		Updated: time.Now(),
	})
	if changed {
		go t.writeSync(key, rev, value) // synchronize downstream tubes
	}
	return
}

// writeSync pushes an update to our downstream peering tubes.
func (t *Tube) writeSync(key string, rev Rev, value interface{}) {
	var wg sync.WaitGroup
	for _, downAvatar := range t.folk.Opened() {
		ydown := YTube{
			tissue.FolkAvatar(downAvatar),
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ydown.Write(key, rev, value)
		}()
	}
	wg.Wait()
}

// BulkWrite writes the given set of records to our local view.
// Resulting changes are pushed to downstream peers as necessary and in bulk.
func (t *Tube) BulkWrite(bulk []*Record) {
	// log.Printf("tube bulk writing")
	// defer func() {
	// 	log.Printf("tube bulk written to\n%s", t.Dump())
	// }()
	if len(bulk) == 0 {
		return
	}
	t.Lock()
	defer t.Unlock()
	changed := make([]*Record, 0, len(bulk))
	for _, r := range bulk {
		if t.view.Update(r) {
			changed = append(changed, r)
		}
	}
	go t.bulkWriteSync(changed)
}

// bulkWriteSync pushes the changed records to our downstream peers.
func (t *Tube) bulkWriteSync(changed []*Record) {
	// Records exchanged within and across tubes are immutable, so no lock is necessary
	var wg sync.WaitGroup
	for _, downAvatar := range t.folk.Opened() {
		ydown := YTube{tissue.FolkAvatar(downAvatar)}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ydown.BulkWrite(changed)
		}()
	}
	wg.Wait()
}

// Forget…
func (t *Tube) Forget(key string, notAfterRev Rev, notAfterUpdated time.Time) bool {
	// log.Printf("tube forgetting %s not after rev %v and not after updated %v", key, notAfterRev, notAfterUpdated)
	// defer func() {
	// 	log.Printf("tube forgot\n%s\n", t.Dump())
	// }()
	t.Lock()
	defer t.Unlock()
	return t.view.Forget(key, notAfterRev, notAfterUpdated)
}

// Scrub…
func (t *Tube) Scrub(key string, notAfterRev Rev, notAfterUpdated time.Time) {
	t.Lock()
	defer t.Unlock()
	if t.view.Forget(key, notAfterRev, notAfterUpdated) {
		go t.scrubSync(key, notAfterRev, notAfterUpdated)
	}
}

// scrub pushes a notification to scrub a record to our downstream peers.
func (t *Tube) scrubSync(key string, notAfterRev Rev, notAfterUpdated time.Time) {
	var wg sync.WaitGroup
	for _, downAvatar := range t.folk.Opened() {
		ydown := YTube{
			tissue.FolkAvatar(downAvatar),
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ydown.Scrub(key, notAfterRev, notAfterUpdated)
		}()
	}
	wg.Wait()
}

// Lookup…
func (t *Tube) Lookup(key string) *Record {
	t.Lock()
	defer t.Unlock()
	return t.view.Lookup(key)
}
