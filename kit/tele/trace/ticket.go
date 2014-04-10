// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package trace

import (
	"fmt"
	"sync"
)

var (
	tlk sync.Mutex
	tkt uint32
)

func TakeTicket() string {
	tlk.Lock()
	defer tlk.Unlock()
	tkt++
	return fmt.Sprintf("TKT:%-2d", tkt)
}
