// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
)

func init() {
	gob.Register(Run{})
}

// Run parameterizes a container execution.
type Run struct {
	Image string
	Memory int64
	CpuShares int64
	Lxc []string
	Volume []string
	Dir string
	Entry string
	Env []string
	Path string
	Args []string
	Scrub bool
}

func ParseRun(src string) (*Run, error) {
	x := &Run{}
	if err := json.Unmarshal([]byte(src), x); err != nil {
		return nil, err
	}
	return x, nil
}

func (x *Run) Arg(name string) []string {
	var r = []string{"run"}
	r = append(r, "--name", name) // name
	if x.CpuShares > 0 {
		r = append(r, "-c", fmt.Sprintf("%d", x.CpuShares))
	}
	if x.Memory > 0 {
		r = append(r, "-m", fmt.Sprintf("%d", x.Memory))
	}
	for _, l := range x.Lxc {
		r = append(r, "--lxc-conf", fmt.Sprintf("%s", l))
	}
	for _, v := range x.Volume {
		r = append(r, "--volume", fmt.Sprintf("%s", v))
	}
	for _, e := range x.Env {
		r = append(r, "--env", fmt.Sprintf("%s", e))
	}
	if x.Dir != "" {
		r = append(r, "--workdir", fmt.Sprintf("%s", x.Dir))
	}
	if x.Entry != "" {
		r = append(r, "--entrypoint", fmt.Sprintf("%s", x.Entry))
	}
	r = append(r, x.Image) // image
	if x.Path != "" {
		r = append(r, x.Path) // command path
	}
	for _, a := range x.Args {
		r = append(r, a)
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
