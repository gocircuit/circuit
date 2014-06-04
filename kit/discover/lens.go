// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package discover

import (
	"sync"

	"github.com/gocircuit/circuit/kit/xor"
)

type family struct {
	cap int // capacity of proximity buckets
	pov xor.Key // point of view
	sync.Mutex
	mem map[int]map[xor.Key]struct{} // proximity => key => {}
}

func newFamily(pov xor.Key, k int) *family {
	if k < 1 {
		panic(0)
	}
	return &family{
		cap: k,
		pov: pov,
		mem: make(map[int]map[xor.Key]struct{}),
	}
}

func (f *family) Clear() {
	f.Lock()
	defer f.Unlock()
	f.mem = make(map[int]map[xor.Key]struct{})
}

func (f *family) Remember(key xor.Key) bool {
	f.Lock()
	defer f.Unlock()
	p := xor.Proximity(f.pov, key)
	s, ok := f.mem[p]
	if !ok {
		s = make(map[xor.Key]struct{})
		f.mem[p] = s
	}
	if len(s) >= f.cap {
		return false
	}
	if _, ok = s[key]; ok { // already accepted
		return false
	}
	s[key] = struct{}{}
	return true
}
