// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Run parameterizes a container execution.
type Run struct {
	Image string `json:"img"`
	Memory int64 `json:"mem"`
	CpuShares int64 `json:"cpu"`
	Lxc []string `json:"lxc"`
	Volume []string `json:"vol"`
	Dir string `json:"dir"`
	Entry string `json:"entry"`
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

func (x *Run) Arg() []string {
	var r = []string{"run"}
	r = append(r, "--net=bridge") // network
	r = append(r, fmt.Sprintf("-c=%d", x.CpuShares)) // cpu shares
	if x.Memory > 0 {
		r = append(r, fmt.Sprintf("-m=%d", x.Memory))
	}
	for _, l := range x.Lxc {
		r = append(r, fmt.Sprintf("--lxc-conf=%s", strconv.QuoteToASCII(l)))
	}
	for _, v := range x.Volume {
		r = append(r, fmt.Sprintf("--volume=%s", strconv.QuoteToASCII(v)))
	}
	for _, e := range x.Env {
		r = append(r, fmt.Sprintf("--env=%s", strconv.QuoteToASCII(e)))
	}
	if x.Dir != "" {
		r = append(r, fmt.Sprintf("--workdir=%s", strconv.QuoteToASCII(x.Dir)))
	}
	if x.Entry != "" {
		r = append(r, fmt.Sprintf("--entrypoint=%s", strconv.QuoteToASCII(x.Entry)))
	}
	r = append(r, x.Image) // image
	if x.Path != "" {
		r = append(r, x.Path) // command path
	}
	for _, a := range x.Args {
		r = append(r, strconv.QuoteToASCII(a))
	}
	return r
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
