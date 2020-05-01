// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tele

import (
	"github.com/hoijui/circuit/pkg/kit/tele/blend"
	"github.com/hoijui/circuit/pkg/kit/tele/codec"
	"github.com/hoijui/circuit/pkg/kit/tele/hmac"
	"github.com/hoijui/circuit/pkg/kit/tele/tcp"
	"github.com/hoijui/circuit/pkg/kit/tele/trace"
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


func NewStructOverTCPWithHMAC(key []byte) *blend.Transport {
	f := trace.NewFrame("tele")
	// Carrier
	x2 := hmac.NewTransport(key)
	// Codec
	x3 := codec.NewTransport(x2, codec.GobCodec{})
	// Blend
	return blend.NewTransport(f.Refine("blend"), x3)
}
