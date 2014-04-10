// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"encoding/gob"
)

func init() {
	gob.Register(&cargo{})
}

type cargo struct {
	Cargo []byte
}

/*
type faithfulReader struct {
	sync.Mutex
	conn *faithful.Conn
	buf  bytes.Buffer
}

func NewFaithfulReader(conn *faithful.Conn) *faithfulReader {
	return &faithfulReader{}
}

func (x *faithfulReader) Read(p []byte) (int, error) {
	x.Lock()
	defer x.Unlock()
	for x.buf.Len() == 0 {
		blob, err := x.conn.Read()
	}
	return x.buf.Read(p)
}
*/
