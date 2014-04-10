// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"reflect"
)

func unflattenValue(v reflect.Value, t reflect.Type) reflect.Value {
	// When t is an Interface, we can't do much, since we don't know the
	// original (unflattened) type of the value placed in v, so we just nop it.
	if t.Kind() == reflect.Interface {
		return v
	}
	// v can be invalid, if it holds the nil value for pointer type
	if !v.IsValid() {
		return v
	}
	// Make sure v is indeed flat
	if v.Kind() == reflect.Ptr {
		panic("unflattening non-flat value")
	}
	// Add a *, one at a time
	for t.Kind() == reflect.Ptr {
		if v.CanAddr() {
			v = v.Addr()
		} else {
			pw := reflect.New(v.Type())
			pw.Elem().Set(v)
			v = pw
		}
		t = t.Elem()
	}
	return v
}

func unflattenSlice(s []interface{}, t []reflect.Type) []interface{} {
	for i, v := range s {
		w := unflattenValue(reflect.ValueOf(v), t[i])
		// If type is *T, v can be invalid (nil) before and after the unflatten call
		if w.IsValid() {
			s[i] = w.Interface()
		}
	}
	return s
}
