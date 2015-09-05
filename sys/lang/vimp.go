// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"reflect"
	"sync"

	"github.com/gocircuit/circuit/use/n"
)

type importGroup struct {
	AllowPP bool   // Allow import of re-exported values (PtrPtr)
	ConnPP  n.Conn // If non-nil, acknowledge receipt of Ptr for each PtrPtr
	sync.Mutex
	Err error
}

func (r *Runtime) importValues(values []interface{}, types []reflect.Type, exporter n.Addr, allowPP bool, connPP n.Conn) ([]interface{}, error) {
	ig := &importGroup{
		AllowPP: allowPP,
		ConnPP:  connPP,
	}
	replacefn := func(src, dst reflect.Value) bool {
		return r.importRewrite(src, dst, exporter, ig)
	}
	rewritten := rewriteInterface(replacefn, values).([]interface{})
	if types == nil {
		return rewritten, ig.Err
	}
	return unflattenSlice(rewritten, types), ig.Err
}

func (r *Runtime) importRewrite(src, dst reflect.Value, exporter n.Addr, ig *importGroup) bool {
	switch v := src.Interface().(type) {

	case *ptrMsg:
		if exporter == nil {
			panic("importing non-perm ptr without exporter")
		}
		imph, err := r.imp.Add(v.ID, v.TypeID, exporter, false)
		if err != nil {
			ig.Lock()
			ig.Err = err
			ig.Unlock()
			return true
		}
		dst.Set(reflect.ValueOf(imph.GetPtr(r)))
		// For each imported handle, wait until it is not needed any more,
		// and notify the exporter.
		imph.ScrubWith(func() {
			// log.Printf("ptr scrubbing %s", imph.Type.Name())
			r.imp.Remove(imph.ID)

			conn, err := r.t.Dial(exporter)
			if err != nil {
				return
			}
			defer conn.Close()
			conn.Write(&dropPtrMsg{imph.ID})
		})
		return true

	case *ptrPtrMsg:
		if exporter == nil {
			panic("importing non-perm ptrptr without exporter")
		}
		if !ig.AllowPP {
			panic("PtrPtr values not allowed in context")
		}
		// Acquire a ptr from the source
		ptr, err := r.callGetPtr(v.ID, v.Src)
		if err != nil {
			ig.Lock()
			ig.Err = err
			ig.Unlock()
			return true
		}
		dst.Set(reflect.ValueOf(ptr))
		if ig.ConnPP != nil {
			// Notify the PtrPtr sender
			if err = ig.ConnPP.Write(&gotPtrMsg{v.ID}); err != nil {
				ig.Lock()
				ig.Err = err
				ig.Unlock()
				return true
			}
		}
		return true

	case *permPtrMsg:
		imph, err := r.imp.Add(v.ID, v.TypeID, exporter, true)
		if err != nil {
			ig.Lock()
			ig.Err = err
			ig.Unlock()
			return true
		}
		dst.Set(reflect.ValueOf(imph.GetPermPtr(r)))
		// For each imported handle, wait until it is not needed any more,
		// and just remove from imports table.
		imph.ScrubWith(func() {
			//log.Printf("perm/ptr scrubbing %s", imph.Type.Name())
			r.imp.Remove(imph.ID)
		})
		return true

	case *permPtrPtrMsg:
		imph, err := r.imp.Add(v.ID, v.TypeID, v.Src, true)
		if err != nil {
			ig.Lock()
			ig.Err = err
			ig.Unlock()
			return true
		}
		dst.Set(reflect.ValueOf(imph.GetPermPtr(r)))
		// For each imported handle, wait until it is not needed any more,
		// and just remove from imports table.
		imph.ScrubWith(func() {
			//log.Printf("perm/ptr/ptr scrubbing %s", imph.Type.Name())
			r.imp.Remove(imph.ID)
		})
		return true
	}
	return false
}
