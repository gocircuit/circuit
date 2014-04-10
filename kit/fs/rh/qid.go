// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package rh

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
)

// Q is a unique ID, reflecting on a file identity captured at a fixed point between updates.
type Q struct {
	ID  uint64
	Ver uint64
}

//
func (q Q) Hash64() uint64 {
	f := fnv.New64a()
	if err := binary.Write(f, binary.BigEndian, q); err != nil {
		panic("u")
	}
	return f.Sum64()
}

//
func (q Q) String() string {
	return fmt.Sprintf("%.16x·%d", q.ID, q.Ver)
}

//
func QBytes(q []byte) uint64 {
	f := fnv.New64a()
	f.Write(q)
	return f.Sum64()
}

func QString(s string) uint64 {
	return QBytes([]byte(s))
}
