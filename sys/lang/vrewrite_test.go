// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	//"fmt"
	"reflect"
	"testing"
)

type testStruct struct {
	a int
	B interface{}
	P interface{}
	Q *testStruct
}

type testReplacement struct{}

func testRewrite(src, dst reflect.Value) bool {
	switch src.Interface().(type) {
	case *_ref:
		dst.Set(reflect.ValueOf(&testReplacement{}))
		return true
	}
	return false
}

func TestRewriteValue(t *testing.T) {
	sv := &testStruct{
		a: 3,
		B: Ref(float64(1.1)),
		P: testStruct{B: int(2)},
		Q: &testStruct{B: int(3)},
	}
	/*xsv :=*/ rewriteInterface(testRewrite, sv)
	//fmt.Printf("%#v\n%#v\n", sv, xsv)
	// XXX: Add test
}
