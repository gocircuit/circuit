// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"encoding/json"
)

type Stat struct {
	Run Run `json:"run"`
	Container string `json:"container"`
	Exit error `json:"exit"`
	Phase string `json:"phase"`
}

func (s Stat) String() string {
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(0)
	}
	return string(b)
}
