// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package docker

import (
	"encoding/json"
	"time"
)

type Stat struct {
	ID string
	Created time.Time
	Path string
	Args []string
	Config
	State
	Image string
	NetworkSettings
	ResolvConfPath string
	HostnamePath string
	HostsPath string
	Name string
	Driver string
	ExecDriver string
	Volumes map[string]string
	VolumesRW map[string]bool
	HostConfig
}

func ParseStat(buf []byte) (s *Stat, err error) {
	s = &Stat{}
	if err = json.Unmarshal(buf, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s Stat) String() string {
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(0)
	}
	return string(b)
}

type Config struct {
	Hostname        string
	Domainname      string
	User            string
	Memory          int64  // Memory limit (in bytes)
	MemorySwap      int64  // Total memory usage (memory + swap); set `-1' to disable swap
	CpuShares       int64  // CPU shares (relative weight vs. other containers)
	Cpuset          string // Cpuset 0-2, 0,1
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	PortSpecs       []string // Deprecated - Can be in the format of 8080/tcp
	ExposedPorts    map[Port]struct{}
	Tty             bool // Attach standard streams to a tty, including stdin if it is not closed.
	OpenStdin       bool // Open stdin
	StdinOnce       bool // If true, close stdin after the 1 attached client disconnects.
	Env             []string
	Cmd             []string
	Image           string // Name of the image as it was passed by the operator (eg. could be symbolic)
	Volumes         map[string]struct{}
	WorkingDir      string
	Entrypoint      []string
	NetworkDisabled bool
	OnBuild         []string
}

type State struct {
	Running    bool
	Paused     bool
	Pid        int
	ExitCode   int
	StartedAt  time.Time
	FinishedAt time.Time
	Ghost bool
}

type NetworkSettings struct {
	IPAddress   string
	IPPrefixLen int
	Gateway     string
	Bridge      string
	PortMapping map[string]PortMapping // Deprecated
	Ports       PortMap
}

type PortMapping map[string]string // Deprecated

type HostConfig struct {
	Binds           []string
	ContainerIDFile string
	LxcConf         []KeyValuePair
	Privileged      bool
	PortBindings    PortMap
	Links           []string
	PublishAllPorts bool
	Dns             []string
	DnsSearch       []string
	VolumesFrom     []string
	NetworkMode     NetworkMode
}

type KeyValuePair struct {
	Key   string
	Value string
}

type NetworkMode string

type PortBinding struct {
	HostIp   string
	HostPort string
}

type Port string

type PortMap map[Port][]PortBinding

type PortSet map[Port]struct{}
