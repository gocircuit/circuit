// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hoijui/circuit/sys/lang/types"
	"github.com/hoijui/circuit/use/circuit"
	"github.com/hoijui/circuit/use/errors"
	"github.com/hoijui/circuit/use/n"
	"github.com/hoijui/circuit/use/worker"
)

func (r *Runtime) Kill(addr n.Addr) error {
	return worker.Kill(addr)
}

// RunInBack can only be invoked inside a serveGo.
// For the user, this means that RunInBack can be called inside functions that
// are invoked via circuit.Spawn
func (r *Runtime) RunInBack(fn func()) {
	r.dwg.Add(1)
	go func() {
		defer r.dwg.Done()
		fn()
	}()
}

func (r *Runtime) serveGo(req *goMsg, conn n.Conn) {
	log.Println("Cross-call start")

	// Go guarantees the defer runs even if panic occurs
	var exit string
	defer func() {
		defer func() {
			recover()
			log.Println("Exit: ", exit)
			os.Exit(0)
		}()
		conn.Close()
		// Potentially unnecessary hack to ensure that last message sent to
		// caller is received before we die
		time.Sleep(time.Second)
	}()

	exit = "lookup"
	t := types.FuncTabl.TypeWithID(req.TypeID)
	if t == nil {
		conn.Write(&returnMsg{Err: NewError("reply: no func type")})
		return
	}
	// No need to acknowledge acquisition of re-exported ptrs since,
	// the caller is waiting for a return message anyway
	exit = "import"
	mainID := t.MainID()
	in, err := r.importValues(req.In, t.Func[mainID].InTypes, conn.Addr(), true, nil)
	if err != nil {
		conn.Write(&returnMsg{Err: err})
		return
	}

	exit = "call"
	reply, err := call(t.Zero(), t, mainID, in)

	if err != nil {
		conn.Write(&returnMsg{Err: err})
		return
	}
	expReply, ptrPtr := r.exportValues(reply, conn.Addr())
	err = conn.Write(&returnMsg{Out: expReply})
	r.readGotPtrPtr(ptrPtr, conn)

	// Wait for any daemonized goroutines to complete
	r.dwg.Wait()
	exit = "ok"
}

func (r *Runtime) Spawn(host worker.Host, anchor []string, fn circuit.Func, in ...interface{}) (retrn []interface{}, addr n.Addr, err error) {

	// Catch all errors
	defer func() {
		if p := recover(); p != nil {
			retrn, addr = nil, nil

			var w bytes.Buffer
			pprof.Lookup("goroutine").WriteTo(&w, 2)
			err = errors.NewError(fmt.Sprintf("spawn panic: %#v\nstack:\n%s", p, string(w.Bytes())))
		}
	}()

	addr, err = worker.Spawn(host, anchor...)
	if err != nil {
		return nil, nil, err
	}

	return r.remoteGo(addr, fn, in...), addr, nil
}

func (r *Runtime) remoteGo(addr n.Addr, ufn circuit.Func, in ...interface{}) []interface{} {
	reply, err := r.tryRemoteGo(addr, ufn, in...)
	if err != nil {
		panic(err)
	}
	return reply
}

// TryGo runs the function ufn on the runtime behind c.
// Any failure to obtain the return values causes a panic.
func (r *Runtime) tryRemoteGo(addr n.Addr, ufn circuit.Func, in ...interface{}) ([]interface{}, error) {
	conn, err := r.t.Dial(addr)
	if err != nil {
		return nil, err
	}
	// Go language spec guarantuees that the defer will run even in the event of panic.
	defer conn.Close()

	expGo, _ := r.exportValues(in, addr)
	t := types.FuncTabl.TypeOf(ufn)
	if t == nil {
		panic(fmt.Sprintf("type '%T' is not a registered worker function type", ufn))
	}
	req := &goMsg{
		// If TypeOf returns nil (causing panic), the user forgot to
		// register the type of ufn
		TypeID: t.ID,
		In:     expGo,
	}
	if err := conn.Write(req); err != nil {
		return nil, NewError("remote write: " + err.Error())
	}
	reply, err := conn.Read()
	if err != nil {
		return nil, NewError("remote read: " + err.Error())
	}
	retrn, ok := reply.(*returnMsg)
	if !ok {
		return nil, NewError("foreign reply")
	}
	if retrn.Err != nil {
		return nil, retrn.Err
	}

	return r.importValues(retrn.Out, t.Func[t.MainID()].OutTypes, addr, true, conn)
}
