// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package anchor

type File struct {
	path   string
	anchor Payload
}

func newFile(h *Handle) *File {
	// File objects do not maintain a handle (strong reference) to the underlying node.
	// They just copy all of the file's immutable data.
	defer h.Close()
	return &File{
		path:   h.Node.Path,
		anchor: h.Node.Anchor.Data,
	}
}

// Path returns the fully-qualified path of this file.
func (f *File) Path() string {
	return f.path
}

// Anchor returns the payload that this file points to.
func (f *File) Anchor() Payload {
	return f.anchor
}
