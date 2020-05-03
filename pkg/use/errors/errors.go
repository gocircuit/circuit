// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package errors

import (
	"encoding/gob"
	"fmt"
	"runtime"
)

func init() {
	gob.Register(&Error{})
}

// NewError creates a simple text-based error that is registered with package
// encoding/gob and therefore can be used in places of error interfaces during
// cross-calls. In contrast, note that due to the rules of gob encoding error objects
// that are not explicitly registered with gob cannot be assigned to error interfaces
// that are to be gob-serialized during a cross-call.
func NewErrorCaller(skip int, fmt_ string, arg_ ...interface{}) error {
	pc, file, line, ok := runtime.Caller(skip + 1)
	var funcname string = "nofunc"
	if ok {
		funcname = runtime.FuncForPC(pc).Name()
	}
	return &Error{
		File: file,
		Line: line,
		Func: funcname,
		Body: fmt.Sprintf(fmt_, arg_...),
	}
}

// NewError
func NewError(fmt_ string, arg_ ...interface{}) error {
	return NewErrorCaller(1, fmt_, arg_...)
}

// NewPanic
func NewPanic(r interface{}) error {
	return NewErrorCaller(1, "panic enclosure (%v)", r)
}

// Pack converts any error into a gob-serializable one that can be used in cross-calls.
func Pack(err error) error {
	if err == nil {
		return nil
	}
	return NewErrorCaller(1, "%s", err.Error())
}

// Error is cross-value that can be used to dynamically wrap a native error value into
// a more informative one for the purposes of passing across workers.
type Error struct {
	File string
	Line int
	Func string
	Body string
}

func (e *Error) Error() string {
	return e.Body
	//return fmt.Sprintf("%s:%d <%s> %s", e.File, e.Line, e.Func, e.Body)
}

func Unpack(x interface{}) error {
	if x == nil {
		return nil
	}
	return x.(error)
}
