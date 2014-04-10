// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package config

import (
	"path"
)

// InstallConfig holds configuration parameters regarding circuit installation on host machines
type InstallConfig struct {
	Dir     string // Root directory of circuit installation on
	LibPath string // Any additions to the library path for execution time
	Worker  string // Desired name for the circuit runtime binary
}

// BinDir returns the binary install directory
func (i *InstallConfig) BinDir() string {
	return path.Join(i.Dir, "bin")
}

// JailDir returns the jail install directory
func (i *InstallConfig) JailDir() string {
	return path.Join(i.Dir, "jail")
}

// VarDir returns the var install directory
func (i *InstallConfig) VarDir() string {
	return path.Join(i.Dir, "var")
}

// BinaryPath returns the absolute path to the worker binary
func (i *InstallConfig) BinaryPath() string {
	return path.Join(i.BinDir(), i.Worker)
}

// ClearHelperPath returns the absolute path to the clear-tool helper binary
func (i *InstallConfig) ClearHelperPath() string {
	return path.Join(i.BinDir(), "4clear-helper")
}
