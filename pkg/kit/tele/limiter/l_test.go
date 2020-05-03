// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package limiter

import (
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	l := New(2)
	for i := 0; i < 9; i++ {
		i_ := i
		l.Go(func() {
			println("{", i_)
			time.Sleep(time.Second)
			println("}", i_)
		})
	}
	l.Wait()
	// TODO: Test that all routines open and close
	println("DONE")
}
