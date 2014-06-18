// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"encoding/json"
	//"strconv"
)

// Run parameterizes a container execution.
type Run struct {
	Image string `json:"image"`
	CpuShares int64 `json:"cpu_shares"`
	Lxc []string `json:"lxc"`
	Memory int64 `json:"memory"`
	Volume []string `json:"volume"`
	Dir string `json:"dir"`
	Env []string `json:"env"`
	Path string `json:"path"`
	Args []string `json:"args"`
}

func ParseRun(src string) (*Run, error) {
	x := &Run{}
	if err := json.Unmarshal([]byte(src), x); err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Run) Volumes() map[string]struct{} {
	m := make(map[string]struct{})
	for _, v := range x.Volume {
		m[v] = struct{}{}
	}
	return m
}

func (x *Run) String() string {
	b, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		panic(0)
	}
	return string(b)
}
