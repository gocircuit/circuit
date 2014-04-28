// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package rh

// NopClunkerFID is an FID that returns without error only from Clunk.
type NopClunkerFID struct {
	ZeroFID
}

func (NopClunkerFID) String() string {
	return ""
}

func (NopClunkerFID) Q() Q {
	return Q{}
}

func (NopClunkerFID) Walk(wname []string) (FID, error) {
	return nil, ErrIO
}

func (NopClunkerFID) Stat() (*Dir, error) {
	return nil, ErrIO
}

func (NopClunkerFID) Open(flag Flag, intr Intr) error {
	return ErrIO
}

func (NopClunkerFID) Clunk() error {
	return nil
}

