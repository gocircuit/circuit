// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

// Package mesos provides the mesos element which serves as a real-time resource megotiation exchange inspired by Apache Mesos.
package mesos

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/hoijui/circuit/use/circuit"
)

type Resource map[string]int // resource name => number of units

type Mesos interface {
	Offer(worker string, rsc Resource)
	Ask(framework string, rsc Resource) (worker string, grant Resource)
}
