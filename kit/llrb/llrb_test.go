// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

import (
	"math"
	"math/rand"
	"testing"
)

func IntLess(p, q interface{}) bool {
	return p.(int) < q.(int)
}

func StringLess(p, q interface{}) bool {
	return p.(string) < q.(string)
}

func TestCases(t *testing.T) {
	tree := New(IntLess)
	tree.ReplaceOrInsert(1)
	tree.ReplaceOrInsert(1)
	if tree.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !tree.Has(1) {
		t.Errorf("expecting to find key=1")
	}

	tree.Delete(1)
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(1) {
		t.Errorf("not expecting to find key=1")
	}

	tree.Delete(1)
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(1) {
		t.Errorf("not expecting to find key=1")
	}
}

func TestReverseInsertOrder(t *testing.T) {
	tree := New(IntLess)
	n := 100
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(n - i)
	}
	c := tree.IterAscend()
	for j, item := 1, <-c; item != nil; j, item = j+1, <-c {
		if item.(int) != j {
			t.Fatalf("bad order")
		}
	}
}

func TestRange(t *testing.T) {
	tree := New(StringLess)
	order := []string{
		"ab", "aba", "abc", "a", "aa", "aaa", "b", "a-", "a!",
	}
	for _, i := range order {
		tree.ReplaceOrInsert(i)
	}
	c := tree.IterRange("ab", "ac")
	k := 0
	for item := <-c; item != nil; item = <-c {
		if k > 3 {
			t.Fatalf("returned more items than expected")
		}
		i1 := string(order[k])
		i2 := item.(string)
		if i1 != i2 {
			t.Errorf("expecting %s, got %s", i1, i2)
		}
		k++
	}
}

func TestRandomInsertOrder(t *testing.T) {
	tree := New(IntLess)
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(perm[i])
	}
	c := tree.IterAscend()
	for j, item := 0, <-c; item != nil; j, item = j+1, <-c {
		if item.(int) != j {
			t.Fatalf("bad order")
		}
	}
}

func TestRandomReplace(t *testing.T) {
	tree := New(IntLess)
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(perm[i])
	}
	perm = rand.Perm(n)
	for i := 0; i < n; i++ {
		if replaced := tree.ReplaceOrInsert(perm[i]); replaced == nil || replaced.(int) != perm[i] {

			t.Errorf("error replacing")
		}
	}
}

func TestRandomInsertSequentialDelete(t *testing.T) {
	tree := New(IntLess)
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(perm[i])
	}
	for i := 0; i < n; i++ {
		tree.Delete(i)
	}
}

func TestRandomInsertDeleteNonExistent(t *testing.T) {
	tree := New(IntLess)
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(perm[i])
	}
	if tree.Delete(200) != nil {
		t.Errorf("deleted non-existent item")
	}
	if tree.Delete(-2) != nil {
		t.Errorf("deleted non-existent item")
	}
	for i := 0; i < n; i++ {
		if u := tree.Delete(i); u == nil || u.(int) != i {
			t.Errorf("delete failed")
		}
	}
	if tree.Delete(200) != nil {
		t.Errorf("deleted non-existent item")
	}
	if tree.Delete(-2) != nil {
		t.Errorf("deleted non-existent item")
	}
}

func TestRandomInsertPartialDeleteOrder(t *testing.T) {
	tree := New(IntLess)
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(perm[i])
	}
	for i := 1; i < n-1; i++ {
		tree.Delete(i)
	}
	c := tree.IterAscend()
	if (<-c).(int) != 0 {
		t.Errorf("expecting 0")
	}
	if (<-c).(int) != n-1 {
		t.Errorf("expecting %d", n-1)
	}
}

func TestRandomInsertStats(t *testing.T) {
	tree := New(IntLess)
	n := 100000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(perm[i])
	}
	avg, _ := tree.HeightStats()
	expAvg := math.Log2(float64(n)) - 1.5
	if math.Abs(avg-expAvg) >= 2.0 {
		t.Errorf("too much deviation from expected average height")
	}
}

func BenchmarkInsert(b *testing.B) {
	tree := New(IntLess)
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(b.N - i)
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	tree := New(IntLess)
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(b.N - i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Delete(i)
	}
}

func BenchmarkDeleteMin(b *testing.B) {
	b.StopTimer()
	tree := New(IntLess)
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(b.N - i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.DeleteMin()
	}
}

func TestInsertNoReplace(t *testing.T) {
	tree := New(IntLess)
	n := 1000
	for q := 0; q < 2; q++ {
		perm := rand.Perm(n)
		for i := 0; i < n; i++ {
			tree.InsertNoReplace(perm[i])
		}
	}
	c := tree.IterAscend()
	for j, item := 0, <-c; item != nil; j, item = j+1, <-c {
		if item.(int) != j/2 {
			t.Fatalf("bad order")
		}
	}
}
