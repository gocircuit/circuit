// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chain

import (
	"errors"
	"fmt"
)

var (
	// ErrMisbehave indicates that the remote endpoint is not behaving to protocol.
	ErrMisbehave = errors.New("misbehave")

	errDup = errors.New("duplicate")
)

// IsStitch returns true if err is a stitching error.
func IsStitch(err error) *ConnWriter {
	if err == nil {
		return nil
	}
	es, ok := err.(*ErrStitch)
	if !ok {
		return nil
	}
	return es.Writer
}

type ErrStitch struct {
	SeqNo  SeqNo
	Writer *ConnWriter
}

func (es *ErrStitch) Error() string {
	return fmt.Sprintf("stitch #%d", es.SeqNo)
}
