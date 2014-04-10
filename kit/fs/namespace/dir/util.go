// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package dir

import (
	"reflect"
	"unicode"

	"github.com/gocircuit/circuit/kit/lang"
)

func valueHash(recv reflect.Value) uint64 {
	return uint64(lang.ComputeReceiverID(recv.Interface()))
}

func interfaceHash(recv interface{}) uint64 {
	return valueHash(reflect.ValueOf(recv))
}

func filename(n string) string {
	var m = []byte(n)
	m[0] = byte(unicode.ToLower(rune(m[0])))
	return string(m)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
