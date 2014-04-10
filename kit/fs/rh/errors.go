// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package rh provides a the interface for the Resource Hierarchy file system
package rh

import (
	"encoding/gob"
	"log"

	"github.com/gocircuit/circuit/use/errors"
)

// Error
type Error interface {
	IsEqual(error) bool
	__RHError()
}

// errorBody
type errorBody string

func (e errorBody) IsEqual(x error) bool {
	f, ok := x.(errorBody)
	if !ok {
		return false
	}
	return string(e) == string(f)
}

func (e errorBody) __RHError() {}

func (e errorBody) Error() string {
	return "rh: " + string(e)
}

var (
	// Hierarchy errors
	ErrExist    = errorBody("exist")     // EBUSY
	ErrNotExist = errorBody("not exist") // ENOENT

	// File descriptor errors
	ErrBusy = errorBody("busy") // EBUSY
	ErrGone = errorBody("gone") // EBADF

	// Security errors
	ErrPerm = errorBody("permission") // EPERM

	// Channel errors
	ErrIO  = errorBody("i/o") // EIO, io.ErrUnexpectedEOF
	ErrEOF = errorBody("eof") // io.EOF (no fuse equivalent?)

	// Protocol errors
	ErrClash = errorBody("clash") // EINVAL invalid arguments or unexpected semantic behavior against the RH API
	ErrIntr  = errorBody("intr")  // EINTR
)

func Naturalize(err error) error {
	if err == nil {
		return nil
	}
	switch err.(*errors.Error).Body {
	case ErrExist.Error():
		return ErrExist
	case ErrNotExist.Error():
		return ErrNotExist
	case ErrBusy.Error():
		return ErrBusy
	case ErrGone.Error():
		return ErrGone
	case ErrPerm.Error():
		return ErrPerm
	case ErrIO.Error():
		return ErrIO
	case ErrEOF.Error():
		return ErrEOF
	case ErrClash.Error():
		return ErrClash
	case ErrIntr.Error():
		return ErrIntr
	}
	log.Panicf("err=%v err=%T err.Error=%v NotExist.Error=%v", err, err, err.Error(), ErrNotExist.Error())
	panic("x")
}

func init() {
	gob.Register(errorBody(""))
}
