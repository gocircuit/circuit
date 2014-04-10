// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package config

import (
	"fmt"
	"testing"
)

const testSrc = `
{
	"Term": {{ "TERM" | env | val }}
}
`

type testData struct {
	Term string
}

func TestParse(t *testing.T) {
	d := &testData{}
	if err := ParseString(d, testSrc); err != nil {
		t.Fatalf("parse string (%s)", err)
	}
	println(fmt.Sprintf("%#v", d))
}
