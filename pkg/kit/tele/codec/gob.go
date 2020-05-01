// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package codec

import (
	"bytes"
	"encoding/gob"
)

// GobCodec
type GobCodec struct{}

func (GobCodec) NewEncoder() Encoder {
	return NewGobEncoder()
}

func (GobCodec) NewDecoder() Decoder {
	return NewGobDecoder()
}

// GobEncoder
type GobEncoder struct {
	w   writer
	enc *gob.Encoder
}

func NewGobEncoder() *GobEncoder {
	g := &GobEncoder{}
	g.w.Clear()
	g.enc = gob.NewEncoder(&g.w)
	return g
}

func (g *GobEncoder) Encode(v interface{}) ([]byte, error) {
	if err := g.enc.Encode(v); err != nil {
		return nil, err
	}
	return g.w.Flush(), nil
}

// GobDecoder
type GobDecoder struct {
	r   reader
	dec *gob.Decoder
}

func NewGobDecoder() *GobDecoder {
	g := &GobDecoder{}
	g.dec = gob.NewDecoder(&g.r)
	return g
}

func (g *GobDecoder) Decode(p []byte, v interface{}) error {
	g.r.Load(p)
	return g.dec.Decode(v)
}

//
type writer struct {
	buf *bytes.Buffer
}

func (w *writer) Clear() {
	w.buf = new(bytes.Buffer)
}

func (w *writer) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *writer) Flush() []byte {
	defer func() {
		w.buf = new(bytes.Buffer)
	}()
	return w.buf.Bytes()
}

//
type reader struct {
	buf *bytes.Reader
}

func (r *reader) Load(p []byte) {
	r.buf = bytes.NewReader(p)
}

func (r *reader) Read(p []byte) (int, error) {
	return r.buf.Read(p)
}
