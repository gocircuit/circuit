// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package blend

import (
	"encoding/gob"
)

type (
	ConnID uint32
	SeqNo  uint32
)

type PayloadMsg struct {
	SeqNo   SeqNo
	Payload interface{} // User-supplied type that can be coded by the underlying codec
}

type AbortMsg struct {
	Err error
}

type Msg struct {
	ConnID ConnID
	Demux  interface{} // PayloadMsg or AbortMsg
}

func init() {
	gob.Register(&PayloadMsg{})
	gob.Register(&AbortMsg{})
}
