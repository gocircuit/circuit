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

type Init struct {
	Env  []string `json:"env"`
	Path string   `json:"path"`
	Args []string `json:"args"`
}

func ParseInit(src string) (*Init, error) {
	x := &Init{}
	if err := json.Unmarshal([]byte(src), x); err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Init) String() string {
	b, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		panic(0)
	}
	return string(b)
}
