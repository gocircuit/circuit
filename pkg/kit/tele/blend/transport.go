// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package blend

import (
	"net"

	"github.com/hoijui/circuit/pkg/kit/tele/codec"
	"github.com/hoijui/circuit/pkg/kit/tele/trace"
)

type Transport struct {
	frame trace.Frame
	sub   *codec.Transport
	*Dialer
}

func NewTransport(frame trace.Frame, sub *codec.Transport) *Transport {
	t := &Transport{
		frame:  frame,
		sub:    sub,
		Dialer: NewDialer(frame.Refine("dialer"), sub),
	}
	frame.Bind(t)
	return t
}

func (t *Transport) Listen(addr net.Addr) *Listener {
	return NewListener(t.frame.Refine("listener"), t.sub.Listen(addr))
}
