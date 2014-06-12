// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package assemble

import (
	"encoding/json"
)

type TraceMsg struct {
	Origin string // "server" or "client"
	Addr string // address of origin
}

func (m *TraceMsg) Encode() []byte {
	buf, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return buf
}

func Decode(p []byte) (*TraceMsg, error) {
	m := &TraceMsg{}
	if err := json.Unmarshal(p, m); err != nil {
		return nil, err
	}
	return m, nil
}
