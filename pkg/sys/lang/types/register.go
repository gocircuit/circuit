// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package types

import (
	"strconv"
)

var (
	ValueTabl *TypeTabl = makeTypeTabl() // Type table for values
	FuncTabl  *TypeTabl = makeTypeTabl() // Type table for functions
)

// RegisterValue registers the type of x with the type table.
// Types need to be registered before values can be imported.
func RegisterValue(r interface{}) {
	ValueTabl.Add(makeType(r))
}

func LookupValue(r interface{}) string {
	return ValueTabl.TypeOf(r).Name()
}

// RegisterFunc ...
func RegisterFunc(fn interface{}) {
	t := makeType(fn)
	if len(t.Func) != 1 {
		panic("fn type must have exactly one method: " + strconv.Itoa(len(t.Func)))
	}
	FuncTabl.Add(t)
}
