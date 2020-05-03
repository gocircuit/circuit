// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"fmt"
	"github.com/hoijui/circuit/pkg/use/circuit"
	"testing"
)

func TestExportImport(t *testing.T) {
	r := New(NewSandbox())
	x := r.Export(1, 2, map[string]circuit.X{"a": nil, "b": nil})
	fmt.Printf("x=%#v\n", x)
	v, s, err := r.Import(x)
	if err != nil {
		t.Errorf("import (%s)", err)
	}
	fmt.Printf("v=%#v, #s=%d\n", v, len(s))
}
