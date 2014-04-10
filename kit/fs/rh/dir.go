// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package rh

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Dir
type Dir struct {
	Q      Q
	Mode   Mode      // Mode
	Perm   Perm      // Perm
	Atime  time.Time // Last access time
	Mtime  time.Time // Last modification time
	Name   string    // File name (excluding directory path to file); must be / if the file is the root directory of the server
	Length int64     // Length
	Uid    string    // owner name
	Gid    string    // group name
	Aux    DirAux    // FID-specific info
}

//
type DirAux interface {
	String() string
}

func (d *Dir) IsDir() bool {
	return d.Mode.Attr == ModeDir
}

const timeFmt = "2006-01-02/15:04:05"

func (d *Dir) String() string {
	return fmt.Sprintf("{name=%s uid=%s gid=%s q=%v m=%s p=%#o at=%d mt=%s aux=%v l=%d}",
		d.Name, d.Uid, d.Gid, d.Q, d.Mode, d.Perm, d.Atime, d.Mtime.Format(timeFmt), d.Aux != nil, d.Length)
}

// Wdir
type Wdir struct {
	Perm   *Perm
	Mtime  *time.Time
	Length *int64
	Gid    string
	Aux    DirAux
}

func (wd *Wdir) String() string {
	var w bytes.Buffer
	if wd.Perm != nil {
		fmt.Fprintf(&w, "perm=%#o ", wd.Perm)
	}
	if wd.Mtime != nil {
		fmt.Fprintf(&w, "mtm=%s ", wd.Mtime.Format(timeFmt))
	}
	if wd.Length != nil {
		fmt.Fprintf(&w, "len=%d ", *wd.Length)
	}
	if wd.Gid != "" {
		fmt.Fprintf(&w, "gid=%s ", wd.Gid)
	}
	if wd.Aux != nil {
		fmt.Fprintf(&w, "aux=%v ", wd.Aux)
	}
	return w.String()
}

//
func init() {
	gob.Register(Dir{})
	gob.Register(Wdir{})
}

func UID() string {
	return strconv.Itoa(os.Getuid())
}

func GID() string {
	return strconv.Itoa(os.Getgid())
}
