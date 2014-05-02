// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"runtime"

	//"github.com/gocircuit/circuit/kit/element/proc"
	//"github.com/gocircuit/circuit/kit/element/valve"
)

type Anchor interface {
	Walk(walk []string) Anchor
	Content() (Element, map[string]Anchor)
	Set(Element)
}

// handle …
type handle struct {
	*anchor
}

func newHandle(a *anchor) (h *handle) {
	h = &handle{a}
	runtime.SetFinalizer(h, func(h2 *handle) { h2.anchor.recycle() })
	return h
}
