// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tele

import (
	"github.com/gocircuit/circuit/kit/tele/blend"
	"github.com/gocircuit/circuit/kit/tele/codec"
	"github.com/gocircuit/circuit/kit/tele/tcp"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

func NewStructOverTCP() *blend.Transport {
	f := trace.NewFrame("tele")
	// Carrier
	x2 := tcp.CodecTransport
	// Codec
	x3 := codec.NewTransport(x2, codec.GobCodec{})
	// Blend
	return blend.NewTransport(f.Refine("blend"), x3)
}
