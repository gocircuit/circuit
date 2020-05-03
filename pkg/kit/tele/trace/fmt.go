// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package trace

import (
	"fmt"
)

type Op int

const (
	READ = Op(iota)
	WRITE
)

func PrintOp(err error, proto string, op Op, msg fmt.Stringer) string {
	var t string
	switch op {
	case READ:
		t = "READ "
	case WRITE:
		t = "WROTE"
	default:
		t = "UKNWN"
	}
	var e string
	if err != nil {
		e = fmt.Sprintf("ERROR(%s)", err)
	}
	return fmt.Sprintf("EVE %5s %5s %s %s", proto, t, msg, e)
}

func DeferPrintOp(frame Frame, err *error, proto string, op Op, msg fmt.Stringer) {
	frame.Println(PrintOp(*err, proto, op, msg))
}
