// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"runtime"
	"sync"
)

type handle struct {
	a *anchor
}

func newHandle(a *anchor) (h *handle) {
	h = &handle{a}
	runtime.SetFinalizer(h, func(h2 *handle) { h2.a.recycle() })
	return h
}

func (h *handle) Walk([]string) Anchor {
}

func (h *handle) MakeChan(name string) {
}
