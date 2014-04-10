// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package xor

import (
	"fmt"
	"math/rand"
	"testing"
)

const K = 16

func TestXOR(t *testing.T) {
	m := &Metric{}
	for i := 0; i < K; i++ {
		m.Add(Key(i))
	}
	for piv := 0; piv < K; piv++ {
		nearest := m.Nearest(Key(piv), K/2)
		fmt.Println(Key(piv).ShortString(4))
		for _, n := range nearest {
			fmt.Println(" ", n.Key().ShortString(4))
		}
	}
}

const stressN = 1000000

func TestStress(t *testing.T) {
	m := &Metric{}
	var h []Key
	for i := 0; i < stressN; i++ {
		id := Key(rand.Int63())
		h = append(h, id)
		if _, err := m.Add(id); err != nil {
			t.Errorf("add (%s)", err)
		}
	}
	perm := rand.Perm(len(h))
	for _, j := range perm {
		m.Remove(h[j])
	}
}
