// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package fmt

import (
	"os"
	"testing"
)

func TestDeep(t *testing.T) {
	s := []interface{}{"a", 2, "c"}
	Deep(os.Stdout, s)
}
