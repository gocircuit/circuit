// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package types implements the runtime type system
package types

import (
	"fmt"
	"reflect"
	"sync"
	"unicode"
	"unicode/utf8"
)

// TypeID is a handle for a type
type TypeID uint64

// TypeTabl assigns universal IDs to local, Go runtime types.
// These IDs are purely a function of the type's package path and name, thereby
// enabling interoperability between compatible but different binaries of the
// underlying Go code.
type TypeTabl struct {
	lk sync.Mutex
	im map[TypeID]*TypeChar
	tm map[reflect.Type]*TypeChar
}

func makeTypeTabl() *TypeTabl {
	return &TypeTabl{
		im: make(map[TypeID]*TypeChar),
		tm: make(map[reflect.Type]*TypeChar),
	}
}

func (tt *TypeTabl) Add(t *TypeChar) {
	tt.lk.Lock()
	defer tt.lk.Unlock()

	if _, present := tt.im[t.ID]; present {
		// panic("type id collision")
		return
	}
	tt.im[t.ID] = t
	tt.tm[t.Type] = t
}

func (tt *TypeTabl) TypeWithID(id TypeID) *TypeChar {
	tt.lk.Lock()
	defer tt.lk.Unlock()
	return tt.im[id]
}

func (tt *TypeTabl) TypeOf(x interface{}) *TypeChar {
	tt.lk.Lock()
	defer tt.lk.Unlock()
	return tt.tm[reflect.TypeOf(x)]
}

// TypeChar reflects on the methods of a Go type, and maintains exportable IDs
// for its methods
type TypeChar struct {
	ID   TypeID
	name string
	Type reflect.Type         // Go type of the receiver
	Func map[FuncID]*funcChar // Public methods
	Proc map[string]*funcChar
}

// makeType makes a new type structure for the receiver's value type
func makeType(receiver interface{}) *TypeChar {
	t := reflect.TypeOf(receiver)
	//
	var q []string
	switch t.Kind() {
	case reflect.Ptr:
		q = []string{t.Elem().PkgPath(), t.Elem().Name(), t.String()}
	default:
		q = []string{t.PkgPath(), t.Name(), t.String()}
	}
	k := &TypeChar{
		// XXX: According to Go's documentation, Type.String alone is not
		// guaranteed to be a unique type identifier. We bet on including
		// PkgPath and Name.  While better, this is still not unique. Pointer
		// types, *T, for example have PkgName equals "", Name equals "" and
		// String equals the non-qualified name of the containing package plus
		// the type name. Thus pointers to identically-named types in different
		// packages identically named directories will collide.
		ID:   TypeID(sliceStringID64(q)),
		name: fmt.Sprintf("%s·%s(%s)", q[0], q[1], q[2]),
		Type: t,
		Func: make(map[FuncID]*funcChar),
		Proc: make(map[string]*funcChar),
	}

	// Reflect methods
	for j := 0; j < t.NumMethod(); j++ {
		m := t.Method(j)
		if p := makeFunc(m, k); p != nil {
			k.Func[p.ID] = p
			k.Proc[m.Name] = p
		}
	}

	return k
}

// Name returns the canonical name of this circuit type
func (t *TypeChar) Name() string {
	return t.name
}

func (t *TypeChar) FuncWithID(id FuncID) *funcChar {
	return t.Func[id]
}

// Zero returns a new zero value of the underlying type
func (t *TypeChar) Zero() reflect.Value {
	return reflect.Zero(t.Type)
}

// New returns a Value representing a pointer to a new zero value for the underlying type.
func (t *TypeChar) New() reflect.Value {
	return reflect.New(t.Type)
}

// MainID returns the ID of the first method in Func.
// It is useful for types that have exactly one method.
func (t *TypeChar) MainID() FuncID {
	for id, _ := range t.Func {
		return id
	}
	panic("no func")
}

// FuncID uniquely represents a method in the context of an interface type
type FuncID int32

type funcChar struct {
	// ID is a collision resistent hash of the method's signature, which
	// includes the method name, its arguments and its replies.
	ID       FuncID
	Method   reflect.Method
	InTypes  []reflect.Type
	OutTypes []reflect.Type
}

func makeFunc(m reflect.Method, parent *TypeChar) *funcChar {
	if m.PkgPath != "" {
		// This is an unexported method
		return nil
	}
	p := &funcChar{Method: m}
	t := m.Type

	// ID computation
	var sign []string
	sign = append(sign, p.Method.Name)

	// Reflect arguments
	// Note that the 0-th argument is the receiver value
	for i := 1; i < t.NumIn(); i++ {
		at := t.In(i)
		if !isExportedOrBuiltinType(at) {
			return nil
		}
		p.InTypes = append(p.InTypes, at)
		sign = append(sign, at.Name())
		gobFlattenRegister(at)
	}

	// End of argument types, beginning of reply types
	sign = append(sign, "—")

	// Reflect return values
	for j := 0; j < t.NumOut(); j++ {
		rt := t.Out(j)
		if !isExportedOrBuiltinType(rt) {
			return nil
		}
		p.OutTypes = append(p.OutTypes, rt)
		sign = append(sign, rt.Name())
		gobFlattenRegister(rt)
	}

	// Precompute ID
	p.ID = FuncID(sliceStringID32(sign))

	return p
}

// Is this an exported — upper case — name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}
