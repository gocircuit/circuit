// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"encoding/json"
)

// Cmd …
type Cmd struct {
	Env   []string `json:"env"`
	Dir   string   `json:"dir"`
	Path  string   `json:"path"`
	Args  []string `json:"args"`
	Scrub bool     `json:"scrub"`
}

func ParseCmd(src string) (*Cmd, error) {
	x := &Cmd{}
	if err := json.Unmarshal([]byte(src), x); err != nil {
		return nil, err
	}
	return x, nil
}

func (x Cmd) String() string {
	b, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		panic(0)
	}
	return string(b)
}
