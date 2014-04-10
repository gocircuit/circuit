// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package config provides access to the circuit configuration of this worker process
package config

import (
	"github.com/gocircuit/circuit/kit/config"
	"github.com/gocircuit/circuit/use/n"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

// Role determines the context within which this executable was invoked
const (
	Main       = "main"
	Daemonizer = "daemonizer"
	Worker     = "worker"
)

var Role string

// CIRCUIT_ROLE names the environment variable that determines the role of this invokation
const RoleEnv = "CIRCUIT_ROLE"

const XXX_shell_barrier = "¢¢¢\n"

// init determines in what context we are being run and reads the configurations accordingly
func init() {
	// println("+ circuit/load/config ……")
	// defer println("+ circuit/load/config ok")

	// Seed random number generation
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Llongfile)

	DefaultSpark = &SparkConfig{
		ID:       n.ChooseWorkerID(),
		BindAddr: "",         // Don't accept incoming circuit calls from other workers
		Host:     "",         // "
		Anchor:   []string{}, // Don't register within the anchor file system
	}

	Config = &WorkerConfig{}
	Role = os.Getenv(RoleEnv)
	if Role == "" {
		Role = Main
	}
	switch Role {
	case Main:
		readAsMain()
	case Daemonizer:
		fmt.Print(XXX_shell_barrier)
		readAsDaemonizerOrWorker()
	case Worker:
		readAsDaemonizerOrWorker()
	default:
		log.Printf("Circuit role '%s' not recognized\n", Role)
		os.Exit(1)
	}
	if Config.Spark == nil {
		Config.Spark = DefaultSpark
	}
}

func readAsMain() {
	// If CIR is set, it points to a single file that contains all three configuration structures in JSON format.
	cir := os.Getenv("CIR")
	if cir == "" {
		println("CIR environment variable is empty")
		os.Exit(1)
	}
	file, err := os.Open(cir)
	if err != nil {
		log.Printf("Problem opening all-in-one config file (%s)", err)
		os.Exit(1)
	}
	defer file.Close()
	parseBag(file)
}

func readAsDaemonizerOrWorker() {
	parseBag(os.Stdin)
}

// WorkerConfig captures the configuration parameters of all sub-systems
// Depending on context of execution, some will be nil.
// Zookeeper and Install should always be non-nil.
type WorkerConfig struct {
	Spark  *SparkConfig
	Anchor *AnchorConfig
	Deploy *InstallConfig
	Build  *BuildConfig
}

// Config holds the worker configuration of this process
var Config *WorkerConfig

func parseBag(r io.Reader) {
	Config = &WorkerConfig{}
	if err := config.Parse(Config, r); err != nil {
		log.Printf("Problem parsing config (%s)", err)
		os.Exit(1)
	}
	if Config.Deploy == nil {
		Config.Deploy = &InstallConfig{}
	}
}
