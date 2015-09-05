// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package io facilitates sharing io reader/writer/closers across workers.
package io

import (
	"io"
	"runtime"
	"time"

	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	// If we are passing I/O objects cross-worker, we want to ensure that the GC
	// is activated regularly so that reclaimed I/O objects will close their
	// underlying resources in a timely manner.
	go func() {
		for {
			time.Sleep(5 * time.Second)
			runtime.GC()
		}
	}()
}

// client types

// YReader
type YReader struct {
	circuit.X
}

func (y YReader) Read(p []byte) (n int, err error) {
	r := y.Call("Read", len(p))
	q, err := unpackBytes(r[0]), errors.Unpack(r[1])
	if len(q) > len(p) {
		panic("corrupt i/o server")
	}
	copy(p, q)
	if err != nil && err.Error() == "EOF" {
		err = io.EOF
	}
	return len(q), err
}

// YWriter
type YWriter struct {
	circuit.X
}

func (y YWriter) Write(p []byte) (n int, err error) {
	r := y.Call("Write", p)
	return r[0].(int), errors.Unpack(r[1])
}

// YCloser
type YCloser struct {
	circuit.X
}

func (y YCloser) Close() error {
	return errors.Unpack(y.Call("Close")[0])
}

// YReadCloser
type YReadCloser struct {
	YReader
	YCloser
}

func NewYReader(u interface{}) YReader {
	return YReader{u.(circuit.X)}
}

func NewYReadCloser(u interface{}) *YReadCloser {
	return &YReadCloser{YReader{u.(circuit.X)}, YCloser{u.(circuit.X)}}
}

// YWriteCloser
type YWriteCloser struct {
	YWriter
	YCloser
}

func NewYWriteCloser(u interface{}) *YWriteCloser {
	return &YWriteCloser{YWriter{u.(circuit.X)}, YCloser{u.(circuit.X)}}
}

// YReadWriteCloser
type YReadWriteCloser struct {
	YReader
	YWriter
	YCloser
}

func NewYReadWriteCloser(u interface{}) *YReadWriteCloser {
	return &YReadWriteCloser{YReader{u.(circuit.X)}, YWriter{u.(circuit.X)}, YCloser{u.(circuit.X)}}
}

// YReadWriter
type YReadWriter struct {
	YReader
	YWriter
}

func NewYReadWriter(u interface{}) *YReadWriter {
	return &YReadWriter{YReader{u.(circuit.X)}, YWriter{u.(circuit.X)}}
}

// X-types

// XReader is a cross-worker exportable object that exposes an underlying local io.Reader.
type XReader struct {
	io.Reader
}

func (x XReader) Read(n int) ([]byte, error) {
	p := make([]byte, n)
	m, err := x.Reader.Read(p)
	return p[:m], errors.Pack(err)
}

// XWriter is a cross-worker exportable object that exposes an underlying local io.Writer.
type XWriter struct {
	io.Writer
}

func (x XWriter) Write(p []byte) (int, error) {
	n, err := x.Writer.Write(p)
	return n, errors.Pack(err)
}

// XCloser is a cross-worker exportable object that exposes an underlying local io.Writer.
type XCloser struct {
	io.Closer
}

// NewXCloser attaches a finalizer to the object which calls Close.
// In cases when a cross-interface to this object is lost because of a failed remote worker,
// the attached finalizer will ensure that before we forget this object the channel it
// encloses will be closed.
func NewXCloser(u io.Closer) *XCloser {
	x := &XCloser{u}
	runtime.SetFinalizer(x, func(x *XCloser) {
		x.Closer.Close()
	})
	return x
}

func (x XCloser) Close() error {
	return errors.Pack(x.Closer.Close())
}

// XReadWriteCloser
type XReadWriteCloser struct {
	XReader
	XWriter
	*XCloser
}

func NewXReadWriteCloser(u io.ReadWriteCloser) circuit.X {
	return circuit.Ref(&XReadWriteCloser{XReader{u}, XWriter{u}, NewXCloser(u)})
}

// XReadCloser
type XReadCloser struct {
	XReader
	*XCloser
}

func NewXReader(u io.Reader) circuit.X {
	return circuit.Ref(XReader{u})
}

func NewXReadCloser(u io.ReadCloser) circuit.X {
	return circuit.Ref(&XReadCloser{XReader{u}, NewXCloser(u)})
}

// XWriteCloser
type XWriteCloser struct {
	XWriter
	*XCloser
}

func NewXWriteCloser(u io.WriteCloser) circuit.X {
	return circuit.Ref(&XWriteCloser{XWriter{u}, NewXCloser(u)})
}

// XReadWriter
type XReadWriter struct {
	XReader
	XWriter
}

func NewXReadWriter(u io.ReadWriter) circuit.X {
	return circuit.Ref(&XReadWriter{XReader{u}, XWriter{u}})
}

// Init
func init() {
	//
	//circuit.RegisterValue(XReader{})
	//circuit.RegisterValue(XWriter{})
	//circuit.RegisterValue(&XCloser{})
	//
	circuit.RegisterValue(&XReadCloser{})
	circuit.RegisterValue(&XWriteCloser{})
	circuit.RegisterValue(&XReadWriteCloser{})
	//
	circuit.RegisterValue(&XReadWriter{})
}

// Utils

func unpackBytes(x interface{}) []byte {
	if x == nil {
		return nil
	}
	return x.([]byte)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
