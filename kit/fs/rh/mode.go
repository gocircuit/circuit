// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// File modified from original in code.google.com/p/goplan9

package rh

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

//
type Mode struct {
	Attr     ModeAttr
	IsHidden bool
	IsUnix   bool
	Unix     os.FileMode
}

func ifstr(cond bool, s string) string {
	if cond {
		return s
	}
	return ""
}

func (m Mode) String() string {
	if m.IsUnix {
		return fmt.Sprintf("%s/%s%s", m.Attr, m.Unix, ifstr(m.IsHidden, "√"))
	} else {
		return m.Attr.String()
	}
}

//
type ModeAttr byte

func (ma ModeAttr) String() string {
	switch ma {
	case ModeDir:
		return "dir"
	case ModeRef:
		return "ref"
	case ModeFile:
		return "file"
	case ModeLog:
		return "log"
	case ModeIO:
		return "io"
	case ModeMutex:
		return "mutex"
	case ModeUnknown:
		return "unknown"
	}
	return "??"
}

//
const (
	ModeDir = ModeAttr(iota)
	ModeRef
	ModeFile
	ModeLog
	ModeIO
	ModeMutex
	ModeUnknown
)

func init() {
	gob.Register(Mode{})
}

// []byte, []*Dir, etc.
type Chunk interface{}

// Chunks
type (
	ByteChunk []byte // Chunk type used for all but directory modes
	DirChunk  []*Dir // Directory mode chunk type
)

func (dc DirChunk) String() string {
	var w bytes.Buffer
	for _, d := range dc {
		w.WriteString(d.String())
		w.WriteByte('\n')
	}
	return w.String()
}

func init() {
	gob.Register(*new(ByteChunk))
	gob.Register(*new(DirChunk))
}
