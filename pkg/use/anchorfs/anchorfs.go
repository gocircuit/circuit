// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package anchorfs exposes the programming interface for accessing the anchor file system
package anchorfs

import (
	"path"
	"strings"
	"time"

	"github.com/hoijui/circuit/pkg/use/errors"
	"github.com/hoijui/circuit/pkg/use/n"
)

var (
	ErrName     = errors.NewError("anchor name")
	ErrNotFound = errors.NewError("not found")
)

// System represents an anchor file system
type System interface {
	OpenFile(string) (File, error)
	OpenDir(string) (Dir, error)

	// Created returns the anchors created by this worker
	Created() []string
}

// Rev is the sequentially-increasing revision number of a file system object
type Rev int32

// Dir is the interface for a directory of workers in the anchor file system
type Dir interface {

	// Path returns the fully-qualified name of the directory
	Path() string

	// List returns the list of file and directory names within this directory alongside with the respective revision.
	List() (rev Rev, files, dirs []string)

	// Change blocks until the contents of this directory changes relative to its contents at revision sinceRev.
	// It then returns the new revision number and contents.
	Change(sinceRev Rev) (rev Rev, files, dirs []string)

	// ChangeExpire is similar to Change, except it timeouts if a change does not occur within an expire interval.
	ChangeExpire(sinceRev Rev, expire time.Duration) (rev Rev, files, dir []string, err error)

	// OpenFile opens the file, registered by the given worker ID, if it exists
	OpenFile(string) (File, error)

	// OpenDir opens a subdirectory
	OpenDir(string) (Dir, error)
}

// File is the interface of an anchor file system file
type File interface {

	// Path returns the fully-qualified name of the file
	Path() string

	// Anchor returns the worker address of the worker who created this file
	Anchor() n.Addr
}

// SanitizeDir ensures that anchor is a valid directory path in the fs
// and returns its parts and a normalized absolute path string.
func SanitizeDir(anchor string) (parts []string, full string, err error) {
	anchor = path.Clean(anchor)
	if len(anchor) == 0 || anchor[0] != '/' {
		return nil, "", ErrName
	}
	parts = strings.Split(anchor[1:], "/")
	for _, part := range parts {
		if _, err = n.ParseWorkerID(part); err == nil {
			// Directory name elements cannot look like worker IDs
			return nil, "", ErrName
		}
	}
	return parts, "/" + path.Join(parts...), nil
}
