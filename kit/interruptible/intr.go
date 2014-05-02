// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package interruptible

import (
	"errors"
)

type Abort chan<- struct{}

func (a Abort) Abort() {
	close(a)
}

type Intr <-chan struct{}

var ErrIntr = errors.New("interrupted")
