// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package iomisc

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

// PrefixReader returns a reader that reads from r, but only in chunks of
// entire lines and returns the lines on the other side prefixed by prefix.
func PrefixReader(p string, r io.Reader) io.Reader {
	return &prefixReader{Reader: bufio.NewReader(r), prefix: p}
}

type prefixReader struct {
	sync.Mutex
	*bufio.Reader
	bytes.Buffer
	prefix string
}

func (r *prefixReader) Read(p []byte) (int, error) {
	r.Lock()
	defer r.Unlock()
	for {
		if r.Buffer.Len() > 0 {
			return r.Buffer.Read(p)
		}
		line, err := r.Reader.ReadString('\n')
		if line != "" {
			r.Buffer.WriteString(r.prefix)
			r.Buffer.WriteString(line)
			continue
		}
		if err != nil {
			return 0, err
		}
		// Empty line read
	}
	panic("u")
}

// PrefixWriter
func PrefixWriter(prefix string, w io.Writer) io.Writer {
	return &prefixWriter{[]byte(prefix), w, true}
}

type prefixWriter struct {
	prefix []byte
	w      io.Writer
	x      bool
}

func (p *prefixWriter) Write(q []byte) (n int, err error) {
	var j int
	for i, y := range q {
		if p.x {
			p.x = false
			if _, err = p.w.Write(p.prefix); err != nil {
				return j, err
			}
		}
		switch y {
		case '\n':
			n, err = p.w.Write(q[j : i+1])
			if err != nil {
				return j + n, err
			}
			j = i + 1
			p.x = true
		}
	}
	n, err = p.w.Write(q[j:])
	return j + n, err
}
