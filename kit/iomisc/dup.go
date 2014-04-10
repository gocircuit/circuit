package iomisc

import (
	"bytes"
	"io"
	"sync"
)

// Dup splits the reader src into two identical readers r1 and r2 that can be used out of sync.
func Dup(u io.Reader) (r1, r2 io.Reader) {
	d := &dupReader{u: u}
	return d, (*dupReader2)(d)
}

type dupReader struct {
	sync.Mutex
	u  io.Reader
	b1 bytes.Buffer
	b2 bytes.Buffer
}

type dupReader2 dupReader

func (d *dupReader) Read(p []byte) (n int, err error) {
	d.Lock()
	defer d.Unlock()
	l := d.b1.Len()
	if l < len(p) {
		n, err = d.u.Read(p[l:])
		if n > 0 {
			d.b1.Write(p[l : l+n])
			d.b2.Write(p[l : l+n])
		}
	}
	n, _ = d.b1.Read(p)
	return n, err
}

func (d *dupReader2) Read(p []byte) (n int, err error) {
	d.Lock()
	defer d.Unlock()
	l := d.b2.Len()
	if l < len(p) {
		n, err = d.u.Read(p[l:])
		if n > 0 {
			d.b1.Write(p[l : l+n])
			d.b2.Write(p[l : l+n])
		}
	}
	n, _ = d.b2.Read(p)
	return n, err
}
