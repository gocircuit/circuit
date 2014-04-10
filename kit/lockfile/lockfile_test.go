// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lockfile

import (
	"testing"
)

func TestLockFile(t *testing.T) {
	const name = "/tmp/test.lock"
	lock, err := Create(name)
	if err != nil {
		t.Fatalf("create lock (%s)", err)
	}

	if _, err := Create(name); err == nil {
		t.Errorf("re-create lock should not succceed", err)
	}

	if err = lock.Release(); err != nil {
		t.Fatalf("release lock (%s)", err)
	}
}
