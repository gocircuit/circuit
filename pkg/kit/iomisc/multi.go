// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package iomisc

import (
	"io"
)

type multiWriter struct {
	writers []io.Writer
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	for i, w := range t.writers {
		if w == nil {
			continue
		}
		n, err = w.Write(p)
		if err != nil {
			t.writers[i] = nil
			continue
		}
		if n != len(p) {
			// Ignore short writes
			// err = io.ErrShortWrite
			continue
		}
	}
	return len(p), nil
}

// MultiWriter creates a writer that duplicates its writes to all the provided
// writers. Unlike io.MultiWriter, this one will process writes even if the
// writers are broken, making sure that writes to this writer never block.
func MultiWriter(writers ...io.Writer) io.Writer {
	return &multiWriter{writers}
}
