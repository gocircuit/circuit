// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// File modified from original in code.google.com/p/goplan9

package rh

import (
	"encoding/gob"
	"fmt"
	"os"
)

//
type Flag struct {
	Attr          FlagAttr
	Create        bool
	RemoveOnClose bool
	CloseOnExec   bool
	Truncate      bool
	IsUnix        bool // Field Unix holds a flag
	Unix          int  // Unix flag
	Deny          bool // If true, accompanying operation should return rh.ErrPerm
}

func (f Flag) String() string {
	var r, c, t, d = "–", "–", "–", "–"
	if f.RemoveOnClose {
		r = "r"
	}
	if f.CloseOnExec {
		c = "c"
	}
	if f.Truncate {
		t = "t"
	}
	if f.Deny {
		d = "ø"
	}
	if f.IsUnix {
		return fmt.Sprintf("%s/%s%s%s%s(%s)", f.Attr, r, c, t, d, unixFlag(f.Unix))
	} else {
		return fmt.Sprintf("%s/%s%s%s%s", f.Attr, r, c, t, d)
	}
}

func init() {
	gob.Register(Flag{})
}

//
type FlagAttr byte

const (
	ReadOnly = FlagAttr(iota)
	WriteOnly
	ReadWrite
	Exec
)

func (a FlagAttr) String() string {
	switch a {
	case ReadOnly:
		return "ro"
	case WriteOnly:
		return "wo"
	case ReadWrite:
		return "rw"
	case Exec:
		return "ex"
	}
	return "??"
}

//
type unixFlag int

func (u unixFlag) String() string {
	return flagString(int(u), unixFlagNames)
}

var unixFlagNames = []flagName{
	{os.O_RDONLY, "ro"},
	{os.O_WRONLY, "wo"},
	{os.O_RDWR, "rw"},
	{os.O_APPEND, "a"},
	{os.O_CREATE, "c"},
	{os.O_EXCL, "excl"},
	{os.O_SYNC, "sync"},
	{os.O_TRUNC, "trunc"},
}

type flagName struct {
	bit  int
	name string
}

func flagString(f int, names []flagName) string {
	var s string

	if f == 0 {
		return "0"
	}

	for _, n := range names {
		if f&n.bit != 0 {
			s += "+" + n.name
			f &^= n.bit
		}
	}
	if f != 0 {
		s += fmt.Sprintf("%+#x", f)
	}
	return s[1:]
}
