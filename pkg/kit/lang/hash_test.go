// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"testing"
)

type X struct{ int }

func TestReceiverID(t *testing.T) {
	var x X
	px := &x
	py := &x
	println(ComputeReceiverID(px).String())
	println(ComputeReceiverID(py).String())
}
