// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package server

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"runtime/pprof"
	"strconv"

	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
)

// DebugDir 
type DebugDir struct {
	rh.FID
	dir *dir.Dir
}

func NewDebugDir() *DebugDir {
	d := &DebugDir{dir: dir.NewDir()}
	d.FID = d.dir.FID()
	d.dir.AddChild("numcpu", newFile(d.Numcpu))
	d.dir.AddChild("numgo", newFile(d.Numgoroutine))
	d.dir.AddChild("mem", newFile(d.Mem))
	d.dir.AddChild("goroutine", newFile(d.Goroutine))
	d.dir.AddChild("block", newFile(d.Block))
	return d
}

func (dbg *DebugDir) Numcpu() string {
	return strconv.Itoa(runtime.NumCPU())
}

func (dbg *DebugDir) Numgoroutine() string {
	return strconv.Itoa(runtime.NumGoroutine())
}

func (dbg *DebugDir) Mem() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return printStruct(m)
}

func (dbg *DebugDir) Goroutine() string {
	var w bytes.Buffer
	pprof.Lookup("goroutine").WriteTo(&w, 1)
	return w.String()
}

func (dbg *DebugDir) Block() string {
	var w bytes.Buffer
	pprof.Lookup("block").WriteTo(&w, 1)
	return w.String()
}

func printStruct(i interface{}) string {
	var v = reflect.ValueOf(i)
	var t = reflect.TypeOf(i)
	var w bytes.Buffer
	for j := 0; j < t.NumField(); j++ {
		var f = t.Field(j)
		switch f.Type.Kind() {
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			//
			fmt.Fprintf(&w, "%s: %v\n", f.Name, v.Field(j).Interface())
		default: //Array, Chan, Func, Interface, Map, Ptr, Slice, String, Struct, UnsafePointer
			fmt.Fprintf(&w, "%s: (skipping)\n", f.Name)
		}
	}
	return w.String()
}
