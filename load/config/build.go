// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package config

// BuildConfig holds configuration parameters for the automated circuit app build system
type BuildConfig struct {
	Binary string // Has no effect. Use InstallConfig.Binary instead.
	Jail   string // Build jail path on build host

	AppRepo   string   // App repo URL
	AppSrc    string   // App GOPATH relative to app repo; or empty string if app repo meant to be cloned inside a GOPATH
	WorkerPkg string   // User program package that should be built as the circuit worker executable
	CmdPkgs   []string // Any additional command packages to build

	GoRepo    string
	RebuildGo bool // Rebuild Go even if a newer version is not available
	Show      bool

	CGO_CFLAGS  string // User-supplied CGO_CFLAGS for the app build
	CGO_LDFLAGS string // User-supplied CGO_LDFLAGS for the app build
	LDFLAGS     string // go build -ldflags='…'

	CircuitRepo string
	CircuitSrc  string

	Host       string // Host where build takes place
	PrefixPath string // PATH to pre-pend to default PATH environment on build host
	Tool       string // Build tool path on build host
	ShipDir    string // Local directory where built runtime binary and dynamic libraries will be delivered
}
