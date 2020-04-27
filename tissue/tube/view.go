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

	"github.com/hoijui/circuit/kit/pubsub"
)

// View is a folk data structure that maintains a key-value set sorted by key
type View struct {
	arrive   *pubsub.PubSub
	depart  *pubsub.PubSub
	sync.Mutex
	lkp       map[string]int // record key => record index
	img      []*Record // Current state of the record space known to us
}

// NewView returns a new view object.
func NewView() (v *View) {
	v = &View{
		lkp: make(map[string]int),
	}
	v.arrive, v.depart = pubsub.New("join", v.peek), pubsub.New("leave", nil)
	return
}

// NewArrivals returns a subscription for the stream of arriving peer identities.
func (v *View) NewArrivals() *pubsub.Subscription {
	return v.arrive.Subscribe()
}

// NewDepartures returns a subscription for the stream of departing peer identities.
func (v *View) NewDepartures() *pubsub.Subscription {
	return v.depart.Subscribe()
}

// Dump returns a textual representation of the contents of this view
func (v *View) Dump() string {
	var w bytes.Buffer
	for _, r := range v.Peek() {
		fmt.Fprintf(&w, "%s––(%d)––>%v\n", r.Key, r.Rev, r.Value)
	}
	return w.String()
}

// Peek returns a copy of the current state of the view.
func (v *View) Peek() []*Record {
	v.Lock()
	defer v.Unlock()
	r := make([]*Record, len(v.img)) // Make an external copy since the img changes continuously
	for i, w := range v.img {
		r[i] = w.Clone()
	}
	return r
}

func (v *View) peek() []interface{} {
	v.Lock()
	defer v.Unlock()
	r := make([]interface{}, len(v.img)) // Make an external copy since the img changes continuously
	for i, w := range v.img {
		r[i] = w.Clone()
	}
	return r
}

// Update updates the state of the view at key r.Key.
// The first update to a key must have revision bigger than zero, otherwise it won't take effect.
// The returned value reflects whether an actual change took place in the view.
func (v *View) Update(r *Record) (changed bool) {
	v.Lock()
	defer v.Unlock()
	i, ok := v.lkp[r.Key]
	if ok && r.Rev <= v.img[i].Rev {
		return false
	}
	r = r.Clone()
	r.Updated = time.Now()
	if ok {
		v.img[i] = r
	} else {
		v.img = append(v.img, r)
		v.lkp[r.Key] = len(v.img) - 1
	}
	if !ok { // Report only keys that didn't exist
		v.arrive.Publish(r.Clone())
	}
	return true
}

// Forget removes the record for key from the view, only if the current record has revision
// no greater than notAfterRev and has not been updated after notUpdatedAfter.
//
// The notAfterRev condition is not in effect is notAfterRev is zero.
// Similarly, the notUpdatedAfter condition is not in effect if notUpdatedAfter is zero.
//
// The returned value reflects whether a removal takes place.
//
func (v *View) Forget(key string, notAfterRev Rev, notUpdatedAfter time.Time) bool {
	v.Lock()
	defer v.Unlock()
	// Decide if a record is being forgotten
	i, ok := v.lkp[key]
	if !ok {
		return false
	}
	r := v.img[i]
	if notAfterRev != 0 && r.Rev > notAfterRev {
		return false
	}
	if !notUpdatedAfter.IsZero() && r.Updated.After(notUpdatedAfter) {
		return false
	}
	// Remove from key index
	delete(v.lkp, key)
	// Compactify the record slice
	n := len(v.img)
	forgotten := v.img[i]
	v.img[i] = v.img[n-1]
	v.img = v.img[:n-1]
	if i < len(v.img) {
		v.lkp[v.img[i].Key] = i
	}
	v.depart.Publish(forgotten)
	return true
}

// Lookup returns a copy of the record for the given key, if one is present.
func (v *View) Lookup(key string) *Record {
	v.Lock()
	defer v.Unlock()
	i, present := v.lkp[key]
	if !present {
		return nil
	}
	return v.img[i].Clone()
}
