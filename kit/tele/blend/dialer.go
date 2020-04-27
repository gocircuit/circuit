// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package blend

import (
	"net"

	"github.com/hoijui/circuit/kit/tele/codec"
	"github.com/hoijui/circuit/kit/tele/trace"
)

type Dialer struct {
	frame trace.Frame
	sub   *codec.Transport
}

func NewDialer(frame trace.Frame, sub *codec.Transport) *Dialer {
	d := &Dialer{frame: frame, sub: sub}
	d.frame.Bind(d)
	return d
}

func (d *Dialer) DialSession(addr net.Addr, scrb func()) (*DialSession, error) {
	sub, err := d.sub.Dial(addr)
	if err != nil {
		return nil, err
	}
	return newDialSession(d.frame.Refine("dial"), sub, scrb), nil // codec.Dial always returns instantaneously
}
