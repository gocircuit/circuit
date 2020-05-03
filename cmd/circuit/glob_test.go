// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"log"
	"testing"
)

func TestGlob(t *testing.T) {
	walk, ellipses := parseGlob("/X/hola/petar/...")
	log.Printf("w=%v ell=%v", walk, ellipses)
}
