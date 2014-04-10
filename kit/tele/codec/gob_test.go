// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package codec

import (
	"reflect"
	"testing"
)

const gobN = 7

type testBlob struct {
	A int
	B int
}

func TestGobCodec(t *testing.T) {
	enc := (GobCodec{}).NewEncoder()
	dec := (GobCodec{}).NewDecoder()
	u := &testBlob{A: 1, B: 5}
	for i := 0; i < gobN; i++ {
		chunk, err := enc.Encode(u)
		if err != nil {
			t.Fatalf("encode (%s)", err)
		}
		v := &testBlob{}
		if err = dec.Decode(chunk, v); err != nil {
			t.Fatalf("decode (%s)", err)
		}
		if !reflect.DeepEqual(v, u) {
			t.Fatalf("not equal")
		}
	}
}
