// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"log"
	"reflect"
	"sync"

	"github.com/gocircuit/circuit/use/n"
)

type exportGroup struct {
	sync.Mutex
	PtrPtr []*ptrPtrMsg
}

func (r *Runtime) exportValues(values []interface{}, importer n.Addr) ([]interface{}, []*ptrPtrMsg) {
	eg := &exportGroup{}
	rewriter := func(src, dst reflect.Value) bool {
		return r.exportRewrite(src, dst, importer, eg)
	}
	return rewriteInterface(rewriter, values).([]interface{}), eg.PtrPtr
}

func (r *Runtime) exportRewrite(src, dst reflect.Value, importer n.Addr, eg *exportGroup) bool {
	// Serialize cross-runtime pointers
	switch v := src.Interface().(type) {

	case *_permptr:
		pm := &permPtrPtrMsg{ID: v.impHandle().ID, TypeID: v.impHandle().Type.ID, Src: v.impHandle().Exporter}
		dst.Set(reflect.ValueOf(pm))
		return true

	case *_ptr:
		if importer == nil {
			panic("exporting non-perm ptrptr without importer")
		}
		pm := &ptrPtrMsg{ID: v.impHandle().ID, Src: v.impHandle().Exporter}
		dst.Set(reflect.ValueOf(pm))
		eg.Lock()
		eg.PtrPtr = append(eg.PtrPtr, pm)
		eg.Unlock()
		return true

	case *_ref:
		if importer == nil {
			panic("exporting non-perm ptr without importer")
		}
		dst.Set(reflect.ValueOf(r.exportPtr(v.value, importer)))
		return true

	case *_permref:
		dst.Set(reflect.ValueOf(r.exportPtr(v.value, nil)))
		return true
	}

	return false
}

// exportPtr returns *permPtrMsg if importer is nil, and *ptrMsg otherwise.
func (r *Runtime) exportPtr(v interface{}, importer n.Addr) interface{} {
	// Add exported value to export table
	exph := r.exp.Add(v, importer)

	if importer == nil {
		return &permPtrMsg{ID: exph.ID, TypeID: exph.Type.ID}
	}

	// Monitor the importer for liveness.
	// DropPtr the handles upon importer death.
	r.lk.Lock()
	defer r.lk.Unlock()
	_, ok := r.live[importer.WorkerID()]
	if !ok {
		r.live[importer.WorkerID()] = struct{}{}

		// The anonymous function creates a "lifeline" connection to the worker importing v.
		// When this conncetion is broken, v is released.
		go func() {

			// Defer removal of v's handle from the export table to the end of this function
			defer func() {
				r.lk.Lock()
				delete(r.live, importer.WorkerID())
				r.lk.Unlock()
				// DropPtr/forget all exported handles
				r.exp.RemoveImporter(importer)
			}()

			conn, err := r.t.Dial(importer)
			if err != nil {
				// log.Println("problem dialing lifeline to", importer.String(), err.Error())
				return
			}
			defer conn.Close()

			if conn.Write(&dontReplyMsg{}) != nil {
				log.Println("problem writing on lifeline to", importer.String(), err.Error())
				return
			}
			// Read returns when the remote dies and
			// runs the conn into an error
			conn.Read()
		}()
	}
	return &ptrMsg{ID: exph.ID, TypeID: exph.Type.ID}
}
