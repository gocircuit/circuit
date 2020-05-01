// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"encoding/gob"
	"fmt"

	"github.com/hoijui/circuit/pkg/sys/lang/types"
	"github.com/hoijui/circuit/pkg/use/circuit"
	"github.com/hoijui/circuit/pkg/use/n"
)

func init() {
	gob.Register(&exportedMsg{})
	// Func invokation-style commands
	gob.Register(&goMsg{})
	gob.Register(&callMsg{})
	gob.Register(&dialMsg{})
	gob.Register(&getPtrMsg{})
	gob.Register(&returnMsg{})
	// Value-passing internal commands
	gob.Register(&gotPtrMsg{})
	gob.Register(&dontReplyMsg{})
	gob.Register(&dropPtrMsg{})
	// Value-passing sub-messages
	gob.Register(&ptrMsg{})
	gob.Register(&ptrPtrMsg{})
	gob.Register(&permPtrMsg{})
	gob.Register(&permPtrPtrMsg{})
}

// Top-level messages

type exportedMsg struct {
	Value []interface{}
	Stack string
}

// Execute a method call
type callMsg struct {
	ReceiverID circuit.HandleID
	FuncID     types.FuncID
	In         []interface{}
}

// Fork a go routine
type goMsg struct {
	TypeID types.TypeID
	In     []interface{}
}

type returnMsg struct {
	Out []interface{}
	Err error
}

type getPtrMsg struct {
	ID circuit.HandleID
}

type gotPtrMsg struct {
	ID circuit.HandleID
}

// dontReplyMsg is dropped by the receiver and intentionally never replies to.
// It is used to sense the death of a runtime.
type dontReplyMsg struct{}

// dialMsg requests that the receiver send back a handle to its permanent.
type dialMsg struct {
	Service string
}

// The importer of a handle sends a release request to the exporter to
// notify them that the held object is no longer needed.
// This is part of the cross-runtime garbage collection mechanism.
type dropPtrMsg struct {
	ID circuit.HandleID
}

// ptrMsg carries ...
type ptrMsg struct {
	ID     circuit.HandleID
	TypeID types.TypeID
}

func (msg *ptrMsg) HandleID() circuit.HandleID {
	panic("(ptrMsg) not for use")
}

func (msg *ptrMsg) Addr() n.Addr {
	panic("(ptrMsg) not for use")
}

func (msg *ptrMsg) IsX() {}

func (msg *ptrMsg) Call(proc string, in ...interface{}) []interface{} {
	panic("hack: not meant to be used")
}

func (msg *ptrMsg) String() string {
	return fmt.Sprintf("ptrMsg: id=%v type=%v", msg.ID, msg.TypeID)
}

// ptrPtrMsg carries ...
type ptrPtrMsg struct {
	ID  circuit.HandleID
	Src n.Addr
}

func (msg *ptrPtrMsg) HandleID() circuit.HandleID {
	panic("(ptrPtrMsg) not for use")
}

func (msg *ptrPtrMsg) Addr() n.Addr {
	panic("(ptrPtrMsg) not for use")
}

func (msg *ptrPtrMsg) IsX() {}

func (msg *ptrPtrMsg) Call(proc string, in ...interface{}) []interface{} {
	panic("hack: not meant to be used")
}

func (msg *ptrPtrMsg) String() string {
	return fmt.Sprintf("ptrPtrMsg: id=%v src=%v", msg.ID, msg.Src)
}

// permPtrMsg carries ...
type permPtrMsg struct {
	ID     circuit.HandleID
	TypeID types.TypeID
}

func (msg *permPtrMsg) HandleID() circuit.HandleID {
	panic("(permPtrMsg) not for use")
}

func (msg *permPtrMsg) Addr() n.Addr {
	panic("(permPtrMsg) not for use")
}

func (msg *permPtrMsg) IsX() {}

func (msg *permPtrMsg) IsPermX() {}

func (msg *permPtrMsg) Call(proc string, in ...interface{}) []interface{} {
	panic("hack: not meant to be used")
}

func (msg *permPtrMsg) String() string {
	return fmt.Sprintf("permPtrMsg: id=%v type=%v", msg.ID, msg.TypeID)
}

// permPtrPtrMsg carries a serialized parmenent x-pointer from a sender to a receiver,
// where the value pointed to is not owned by the sender.
type permPtrPtrMsg struct {
	ID     circuit.HandleID
	TypeID types.TypeID
	Src    n.Addr
}

func (msg *permPtrPtrMsg) HandleID() circuit.HandleID {
	panic("(permPtrPtrMsg) not for use")
}

func (msg *permPtrPtrMsg) Addr() n.Addr {
	panic("(permPtrPtrMsg) not for use")
}

func (msg *permPtrPtrMsg) IsX() {}

func (msg *permPtrPtrMsg) IsPermX() {}

func (msg *permPtrPtrMsg) Call(proc string, in ...interface{}) []interface{} {
	panic("hack: not meant to be used")
}

func (msg *permPtrPtrMsg) String() string {
	return fmt.Sprintf("permPtrPtrMsg: id=%v type=%v src=%v", msg.ID, msg.TypeID, msg.Src)
}
