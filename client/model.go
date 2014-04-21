// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package client is development in progress. Do not use, but feel free to look.
package client

import (
	"io"
)

// Process

type Command struct {
	Env  []string `json:"env"`
	Path string   `json:"path"`
	Args []string `json:"args"`
}

// Select

// Clause stands for any of the Clause* types.
type Clause interface{}

// ClauseDefault represnts a default case in a selection.
type ClauseDefault struct{}

// ??
type ClauseSend struct {
	*Chan
}

// ??
type ClauseRecv struct {
	*Chan
}

// ??
type ClauseExit struct {
	*Proc
}
