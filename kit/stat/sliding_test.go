// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package stat

import (
	"testing"
	"time"
)

func TestCarousel(t *testing.T) {
	c := NewSlidingMoment(10, time.Second)

	now := time.Now().UnixNano()
	for i := 0; i < 20; i++ {
		slot := c.Slot(time.Unix(0, now-int64(i)*1e9))
		if slot != nil {
			println("ok")
			slot.Add(1)
		} else {
			println("gh")
		}
	}
	for i := 0; i < 5; i++ {
		slot := c.Slot(time.Unix(0, now+int64(i)*1e9))
		slot.Add(5)
	}
	slots, tlast := c.Slots()
	println("timelast", tlast.UnixNano())
	for _, slot := range slots {
		println("slot", slot.Count(), slot.Max())
	}
	println(c.Weight())
}
