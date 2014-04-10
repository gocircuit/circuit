// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package blend

import (
	"errors"
)

var (
	ErrClash = errors.New("clash")
	ErrGone  = errors.New("gone")
	ErrOff   = errors.New("off")
)
