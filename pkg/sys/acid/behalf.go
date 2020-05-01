// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package acid

import (
	"github.com/hoijui/circuit/pkg/use/circuit"
	"fmt"
	"reflect"
)

type Stringer interface {
	String() string
}

func (s *Acid) OnBehalfCallStringer(service, proc string) (r string) {

	// If anything goes wrong, let's not panic the worker
	defer func() {
		if p := recover(); p != nil {
			r = fmt.Sprintf("Stat likely not supported:\n%#v", p)
		}
	}()

	// Obtain service object in this worker
	srv := circuit.DialSelf(service)
	if srv == nil {
		return "Service not available"
	}

	// Find Stat method in service receiver s
	sv := reflect.ValueOf(srv)
	out := sv.MethodByName(proc).Call(nil)
	if len(out) != 1 {
		return "Service's Stat method returns more than one value"
	}

	return out[0].Interface().(Stringer).String()
}
