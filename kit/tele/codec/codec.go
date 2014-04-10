// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package codec

// Codec
type Codec interface {
	NewEncoder() Encoder
	NewDecoder() Decoder
}

// Encoder
type Encoder interface {
	Encode(interface{}) ([]byte, error)
}

// Decoder
type Decoder interface {
	Decode([]byte, interface{}) error
}

// ChunkCodec
type ChunkCodec struct{}

func (ChunkCodec) NewEncoder() Encoder {
	return ChunkEncoder{}
}

func (ChunkCodec) NewDecoder() Decoder {
	return ChunkDecoder{}
}

// ChunkEncoder
type ChunkEncoder struct{}

func (ChunkEncoder) Encode(v interface{}) ([]byte, error) {
	return v.([]byte), nil
}

// ChunkDecoder
type ChunkDecoder struct{}

func (ChunkDecoder) Decode(r []byte, v interface{}) error {
	v = r
	return nil
}
