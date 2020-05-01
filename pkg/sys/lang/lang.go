// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package lang implements the language runtime
package lang

import (
	"github.com/hoijui/circuit/pkg/kit/lang"
	"github.com/hoijui/circuit/pkg/use/circuit"
	"github.com/hoijui/circuit/pkg/use/n"
)

// _ref wraps a user object, indicating to the runtime that the user has
// elected to send this object as a ptr across runtimes.
type _ref struct {
	value interface{}
}

func (*_ref) Addr() n.Addr {
	return circuit.ServerAddr()
}

func (r *_ref) HandleID() circuit.HandleID {
	panic("handle not assigned yet")
}

func (x *_ref) String() string {
	return "xref://" + lang.ComputeReceiverID(x.value).String()
}

func (*_ref) IsX() {}

func (*_ref) Call(proc string, in ...interface{}) []interface{} {
	panic("call on ref")
}

// _permref
type _permref struct {
	value interface{}
}

func (x *_permref) String() string {
	return "xpermref://" + lang.ComputeReceiverID(x.value).String()
}

func (*_permref) Addr() n.Addr {
	return circuit.ServerAddr()
}

func (pr *_permref) HandleID() circuit.HandleID {
	panic("handle not assigned yet")
}

func (*_permref) IsX() {}

func (*_permref) IsPermX() {}

func (*_permref) Call(proc string, in ...interface{}) []interface{} {
	panic("call on permref")
}

// Ref annotates a user value v, so that if the returned value is consequently
// passed cross-runtime, the runtime will pass v as via a cross-runtime pointer
// rather than by value.
func (*Runtime) Ref(v interface{}) circuit.X {
	if v == nil {
		return nil
	}
	return Ref(v)
}

func Ref(v interface{}) circuit.X {
	if v == nil {
		return nil
	}
	switch v := v.(type) {
	case *_ptr:
		return v
	case *_ref:
		return v
	case *_permptr:
		return v
	case *_permref:
		panic("applying ref on permref")
	}
	return &_ref{v}
}

func (*Runtime) PermRef(v interface{}) circuit.PermX {
	if v == nil {
		return nil
	}
	return PermRef(v)
}

func PermRef(v interface{}) circuit.PermX {
	if v == nil {
		return nil
	}
	switch v := v.(type) {
	case *_ptr:
		panic("permref on ptr")
	case *_ref:
		panic("permref on ref")
	case *_permptr:
		return v
	case *_permref:
		return v
	}
	return &_permref{v}
}

func IsX(v interface{}) bool {
	if v == nil {
		return false
	}
	switch v.(type) {
	case *_ptr, *_ref, *_permptr, *_permref:
		return true
	}
	return false
}
