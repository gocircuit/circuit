// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/gocircuit/circuit/sys/lang/types"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
)

// expTabl issues (and reclaims) universal handles to local values
// and matching local handle structures
type expTabl struct {
	tt      *types.TypeTabl
	lk      sync.Mutex
	id      map[circuit.HandleID]*expHandle
	perm    map[interface{}]*expHandle
	nonperm map[n.WorkerID]map[interface{}]*expHandle // Worker ID -> receiver -> export handle
}

// expHandle holds the underlying local value of an exported handle
type expHandle struct {
	ID       circuit.HandleID
	Importer n.Addr
	Value    reflect.Value // receiver of methods
	Type     *types.TypeChar
}

func (exph *expHandle) String() string {
	return fmt.Sprintf("ID:%s Importer:%#v Type:%s", exph.ID, exph.Importer, exph.Type.Name())
}

func makeExpTabl(tt *types.TypeTabl) *expTabl {
	return &expTabl{
		tt:      tt,
		id:      make(map[circuit.HandleID]*expHandle),
		perm:    make(map[interface{}]*expHandle),
		nonperm: make(map[n.WorkerID]map[interface{}]*expHandle),
	}
}

func (exp *expTabl) Add(receiver interface{}, importer n.Addr) *expHandle {
	if receiver == nil {
		panic("bug: nil receiver in export")
	}

	exp.lk.Lock()
	defer exp.lk.Unlock()

	// Is receiver already exported in the respective perm/nonperm fashion?
	var impHere bool
	var impTabl map[interface{}]*expHandle
	if importer != nil {
		// Non-permanent case
		impTabl, impHere = exp.nonperm[importer.WorkerID()]
		if impHere {
			exph, present := impTabl[receiver]
			if present {
				return exph
			}
		}
	} else {
		// Permanent case
		if exph, present := exp.perm[receiver]; present {
			return exph
		}
	}

	// Build exported handle object
	//fmt.Printf("recv (%#T): %#v\n", receiver, receiver)
	typ := exp.tt.TypeOf(receiver)
	if typ.Type != reflect.TypeOf(receiver) {
		panic("bug: wrong type")
	}
	exph := &expHandle{
		ID:       circuit.ChooseHandleID(),
		Importer: importer,
		Value:    reflect.ValueOf(receiver),
		Type:     typ,
	}

	// Insert in handle map
	if _, present := exp.id[exph.ID]; present {
		panic("handle id collision")
	}
	exp.id[exph.ID] = exph

	// Insert in value map
	if importer != nil {
		// Non-permanent case
		if !impHere {
			impTabl = make(map[interface{}]*expHandle)
			exp.nonperm[importer.WorkerID()] = impTabl
		}
		impTabl[receiver] = exph
	} else {
		// Permanent case
		exp.perm[receiver] = exph
	}

	return exph
}

func (exp *expTabl) Lookup(id circuit.HandleID) *expHandle {
	exp.lk.Lock()
	defer exp.lk.Unlock()

	return exp.id[id]
}

// Remove removes the exported value with handle id from the table, if present.
// If present, a check is performed that importer is the same one, registered
// with the table. If not, an error is returned.
func (exp *expTabl) Remove(id circuit.HandleID, importer n.Addr) {
	if importer == nil {
		panic("cannot remove perm handles from exp")
	}
	exp.lk.Lock()
	defer exp.lk.Unlock()

	exph, present := exp.id[id]
	if !present {
		return
	}
	if importer.WorkerID() != exph.Importer.WorkerID() {
		panic("releasing importer different than original")
	}
	delete(exp.id, id)

	impTabl, present := exp.nonperm[exph.Importer.WorkerID()]
	if !present {
		panic("missing importer map")
	}
	delete(impTabl, exph.Value.Interface())

	if len(impTabl) == 0 {
		delete(exp.nonperm, exph.Importer.WorkerID())
	}
}

func (exp *expTabl) RemoveImporter(importer n.Addr) {
	if importer == nil {
		panic("nil importer")
	}
	exp.lk.Lock()
	defer exp.lk.Unlock()

	impTabl, present := exp.nonperm[importer.WorkerID()]
	if !present {
		return
	}
	delete(exp.nonperm, importer.WorkerID())

	for _, exph := range impTabl {
		delete(exp.id, exph.ID)
	}
	runtime.GC()
}
