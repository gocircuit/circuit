// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/gocircuit/circuit/use/circuit"
)

type Bulletin interface {
	WaitJoin() string
	WaitLeave() string
	IsDone() bool
	Scrub()
	Stat() Stat
	X() circuit.X
}

type bulletin struct {
}

func New() Bulletin {
	??
}
