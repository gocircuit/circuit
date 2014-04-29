// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gocircuit/circuit/kit/iomisc"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

type Error struct {
	Error string `json:"error"`
}

func NewErrorFile() *ErrorFile {
	return &ErrorFile{}
}

type ErrorFile struct {
	sync.Mutex
	err interface{} // err can be any object that can JSON-marshal
}

func (f *ErrorFile) Clear() {
	f.Lock()
	defer f.Unlock()
	f.err = nil
}

func (f *ErrorFile) Set(msg string) {
	f.Lock()
	defer f.Unlock()
	f.err = Error{msg}
}

func (f *ErrorFile) Setf(format string, arg ...interface{}) {
	f.Set(fmt.Sprintf(format, arg...))
}

func (f *ErrorFile) Get() interface{} {
	f.Lock()
	defer f.Unlock()
	return f.err
}

func (f *ErrorFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *ErrorFile) Open(rh.Flag, rh.Intr) (rh.FID, error) {
	return NewOpenReaderFile(iomisc.ReaderNopCloser(bytes.NewBuffer(marshal(f.Get())))), nil
}

func (f *ErrorFile) Remove() error {
	return rh.ErrPerm
}

func marshal(v interface{}) []byte {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return b
}
