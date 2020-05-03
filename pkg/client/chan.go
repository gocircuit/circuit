// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"io"

	"github.com/hoijui/circuit/pkg/element/valve"
)

// Chan provides access to a circuit channel element.
//
// A channel element is semantically identical to a Go channel, with the sole exception
// that the "messages" passed through the channel are pipes that connect the sender
// and the receiver and allow them, once connected, to exchange an arbitrary stream of
// byte data which as a whole counts as one channel message.
//
// All methods panic if the server hosting the channel dies.
type Chan interface {

	// Send blocks until the requested transmission is matched to a receiving call to Recv, or
	// until it can be accommodated in the channel's buffer.
	// It returns a WriteCloser representing a byte pipe to the receiver, or a non-nil error
	// if the channel has already been closed.
	Send() (io.WriteCloser, error)

	// Scrub aborts and abandons the channel. Any buffered send operations are lost.
	Scrub()

	// Close closes the channel, reporting an error only if the channel has already been closed.
	Close() error

	// Recv blocks until it can be matched with a sender.
	// It returns a ReadCloser for the byte pipe from the sender, or a non-nil error if the
	// channel has been closed.
	Recv() (io.ReadCloser, error)

	// Cap reports the capacity of the channel.
	Cap() int

	// Stat returns the current state of the channel.
	Stat() ChanStat
}

// ChanStat describes the state of a channel.
type ChanStat struct {

	// Cap is the channel capacity.
	Cap int

	// Closed is set as soon as Close is called.
	// If Closed is set, the channel might still have messages in its buffer and
	// thus its receive side remains operational.
	Closed bool

	// Aborted is set if the channel has been permanently aborted and is not usable any longer.
	Aborted bool

	// NumSend is the number of completed invocations to Send.
	NumSend int

	// NumRecv is the number of completed invocations to Recv.
	NumRecv int
}

func retypeChanStat(s valve.Stat) ChanStat {
	return ChanStat{
		Cap:     s.Cap,
		Closed:  s.Closed,
		Aborted: s.Aborted,
		NumSend: s.NumSend,
		NumRecv: s.NumRecv,
	}
}

type yvalveChan struct {
	valve.YValve
}

func (y yvalveChan) Stat() ChanStat {
	return retypeChanStat(y.YValve.Stat())
}
