// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/gocircuit/circuit/kit/iomisc"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

type CapFile struct {
	v *Valve
}

func NewCapFile(v *Valve) file.File {
	return &CapFile{v: v}
}

func (f *CapFile) Perm() rh.Perm {
	return 0666 // rw-rw-rw-
}

func (f *CapFile) Open(flag rh.Flag, intr rh.Intr) (rh.FID, error) {
	switch flag.Attr {
	case rh.ReadOnly:
		return file.NewOpenReaderFile(
			iomisc.ReaderNopCloser(
				bytes.NewBufferString(
					strconv.Itoa(f.v.GetCap()) + "\n",
				),
			),
		), nil
	case rh.WriteOnly:
		return file.NewOpenWriterFile(&capWriteFile{v: f.v}), nil
	}
	return nil, rh.ErrPerm
}

func (f *CapFile) Remove() error {
	return rh.ErrPerm
}

type capWriteFile struct {
	v *Valve
	bytes.Buffer
}

func (c *capWriteFile) Close() error {
	c.v.ErrorFile.Clear()
	n, err := strconv.Atoi(strings.TrimSpace(c.Buffer.String()))
	if err != nil || n < 0 {
		c.v.ErrorFile.Set("capacity not a non-negative integer")
		return rh.ErrClash
	}
	return c.v.SetCap(n)
}
