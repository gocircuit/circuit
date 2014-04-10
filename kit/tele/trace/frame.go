// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package trace implements an ad-hoc tracing system
package trace

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"unicode/utf8"
)

var _log = log.New(os.Stderr, "", log.LstdFlags|log.Llongfile)

// Framed is an object that has a Frame
type Framed interface {
	Frame() Frame
}

// Frame …
type Frame interface {
	Refine(sub ...string) Frame
	Bind(interface{})

	Println(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	String() string
	Chain() []string
}

// frame implements Frame
type frame struct {
	ptr uintptr
	chain
}

func NewFrame(s ...string) Frame {
	return &frame{chain: chain(s)}
}

func (f *frame) Refine(sub ...string) Frame {
	c := make(chain, len(f.chain), len(f.chain)+len(sub))
	copy(c, f.chain)
	c = append(c, sub...)
	return &frame{chain: c}
}

func (f *frame) Bind(v interface{}) {
	if f.ptr != 0 {
		panic("duplicate binding")
	}
	f.ptr = reflect.ValueOf(v).Pointer()
}

func justify(s string, l int) string {
	var w bytes.Buffer
	n := utf8.RuneCountInString(s)
	for i := 0; i < max(0, l-n); i++ {
		w.WriteRune('·')
	}
	w.WriteString(s)
	return string(w.Bytes())
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (f *frame) String() string {
	return fmt.Sprintf("(0x%010x) %s", f.ptr, f.chain.String())
}

func (f *frame) Println(v ...interface{}) {
	if os.Getenv("TELE_NOTRACE") != "" {
		return
	}
	_log.Output(2, fmt.Sprintln(append([]interface{}{f.String()}, v...)...))
}

func (f *frame) Print(v ...interface{}) {
	if os.Getenv("TELE_NOTRACE") != "" {
		return
	}
	_log.Output(2, fmt.Sprint(append([]interface{}{f.String()}, v...)...))
}

func (f *frame) Printf(format string, v ...interface{}) {
	if os.Getenv("TELE_NOTRACE") != "" {
		return
	}
	_log.Output(2, fmt.Sprintf("%s "+format, append([]interface{}{f.String()}, v...)...))
}

func (f *frame) Chain() []string {
	return []string(f.chain)
}

// chain is an ordered sequence of strings with a String method
type chain []string

func (c chain) String() string {
	var w bytes.Buffer
	w.WriteString("(")
	for i, s := range c {
		if i > 0 {
			w.WriteString("·")
		}
		w.WriteString(s)
	}
	w.WriteString(")")
	return string(w.Bytes())
}
