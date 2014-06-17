// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"encoding/json"

	"github.com/gocircuit/circuit/element/proc"

	// "github.com/fsouza/go-dockerclient"
)

// Run parameterizes a container execution.
type Run struct {
	Image string `json:"image"`
	CPUShares int `json:"cpu_shares"`
	LXC []string `json:"lxc"`
	Memory string `json:"memory"`
	Volume []string `json:"volume"`
	Dir string `json:"dir"`
	Cmd proc.Cmd `json:"cmd"`
}

func ParseRun(src string) (*Run, error) {
	x := &Run{}
	if err := json.Unmarshal([]byte(src), x); err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Run) String() string {
	b, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		panic(0)
	}
	return string(b)
}
