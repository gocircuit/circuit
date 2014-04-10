// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package faithful

import (
	"io"
)

// BufferRead holds the return values of a call to Buffer.Read
type BufferRead struct {
	Payload interface{}
	SeqNo   SeqNo
	Err     error
}

func NewBufferReadChan(bfr *Buffer) <-chan *BufferRead {
	ch := make(chan *BufferRead)
	go func() {
		defer close(ch)
		for {
			var br BufferRead
			br.Payload, br.SeqNo, br.Err = bfr.Read()
			ch <- &br
			if br.Err == io.ErrUnexpectedEOF {
				return
			}
		}
	}()
	return ch
}
