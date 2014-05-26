// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tube

import (
	"time"
)

// Record is the data unit of the tube system.
type Record struct {
	Key     string
	Rev     Rev
	Value   interface{}
	Updated time.Time
}

// Rev is an increasing revision number of a tube record.
type Rev uint64

func (r *Record) Clone() *Record {
	var p = *r
	return &p
}
