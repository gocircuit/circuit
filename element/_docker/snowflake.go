// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chamber

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocircuit/circuit/kit/docker/docker"
	"github.com/gocircuit/circuit/use/n"
)

// Parent loci spin child loci (docker contains) here.
// The child logic is invoked by circuit/cmd/circuit/load.go, and implemented in child.go

// The Snowflake structure manages a circuit's dockerized subordinate circuits.
type Snowflake struct {
	Config       Config
	ParentDir    string // Local directory holding parent circuit's binaries
	ParentBinary string // Name of parent circuit's binary
	//
	sync.Mutex
	buf struct {
		out, err bytes.Buffer
	}
	chamber map[string]*Chamber // Docker container ID –> chamber
}

type Config struct {
	Port      int    // Child port
	Genus     string // Working directory inside child container
	ParentURL string // URL of this parent worker, for purposes of subordinate kinfolk join
}

func NewSnowflake(config Config) (snowflake *Snowflake, err error) {
	snowflake = &Snowflake{
		Config:       config,
		ParentDir:    path.Dir(os.Args[0]),
		ParentBinary: path.Base(os.Args[0]),
	}
	if snowflake.dc, err = docker.NewDockerCli(nil, &snowflake.buf.out, &snowflake.buf.err, "unix", "/var/run/docker.sock"); err != nil {
		return nil, err
	}
	return
}
