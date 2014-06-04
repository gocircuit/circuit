// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package xor

import "hash/fnv"
import "math/rand"

func ChooseKey() Key {
	return Key(rand.Int63())
}

// ChooseMinK chooses k random keys and returns the one which, if inserted into
// the metric, would result in the shallowest position in the XOR tree.
// In other words, it returns the most balanced choice. We recommend k equals 7.
func (m *Metric) ChooseMinK(k int) Key {
	if m == nil {
		return Key(rand.Int63())
	}
	var min_id Key
	var min_d int = 1000
	for k > 0 {
		// Note: The last bit is not really randomized here
		id := ChooseKey()
		d, err := m.Add(id)
		if err != nil {
			continue
		}
		m.Remove(id)
		if d < min_d {
			min_id = id
			min_d = d
		}
		k--
	}
	return min_id
}

func Combine(keys ...Key) Key {
	r := Key(0)
	for _, k := range keys {
		// This is not very robust. It assumes incoming keys are independently random.
		r ^= k
	}
	return r
}

func uint64bytes(u uint64) []byte {
	q := make([]byte, 8)
	for i, _ := range q {
		q[i] = byte(u >> uint64(i) * 8)
	}
	return q
}

func HashInt64(i int64) Key {
	h := fnv.New64a()
	h.Write(uint64bytes(uint64(i)))
	return Key(h.Sum64())
}

func HashString(s string) Key {
	h := fnv.New64a()
	h.Write([]byte(s))
	return Key(h.Sum64())
}

func HashBytes(b []byte) Key {
	h := fnv.New64a()
	h.Write(b)
	return Key(h.Sum64())
}
