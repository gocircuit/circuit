// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"fmt"
	"github.com/hoijui/circuit/pkg/sys/lang/types"
	"github.com/hoijui/circuit/pkg/use/circuit"
	"runtime"
	"strconv"
	"testing"
	"time"
)

type testBoot struct {
	name string
}

func (x *testBoot) NewGreeter() interface{} {
	return Ref(&testGreeter{})
}

func (x *testBoot) UseGreeter(g circuit.X) string {
	return g.(circuit.X).Call("String")[0].(string)
}

func (x *testBoot) Name() string {
	return x.name
}

func (x *testBoot) ReturnNil() *TestData {
	return nil
}

func (x *testBoot) ReturnNilMap() *TestData {
	m := make(map[string]interface{})
	m["a"] = nil
	return &TestData{M: m}
}

type testGreeter struct{}

const testGreeting = "hey, how are you?"

func (x *testGreeter) String() string {
	return testGreeting
}

func TestX(t *testing.T) {
	l1 := NewSandbox()
	r1 := New(l1)
	r1.Listen("test", &testBoot{"π1"})

	l2 := NewSandbox()
	r2 := New(l2)
	r2.Listen("test", &testBoot{"π2"})

	p2, err := r1.TryDial(l2.Addr(), "test")
	if err != nil {
		t.Fatalf("dial 1->2 (%s)", err)
	}

	p1, err := r2.TryDial(l1.Addr(), "test")
	if err != nil {
		t.Fatalf("dial 2->1 (%s)", err)
	}

	if p1.Call("Name")[0].(string) != "π1" {
		t.Errorf("return val 1")
	}

	if p2.Call("Name")[0].(string) != "π2" {
		t.Errorf("return val 2")
	}
	p2.Call("ReturnNilMap")
}

func TestRe(t *testing.T) {
	r := make([]*Runtime, 3)
	l := make([]n.Transport, len(r))
	for i := 0; i < len(r); i++ {
		l[i] = NewSandbox()
		r[i] = New(l[i])
		r[i].Listen("test", &testBoot{fmt.Sprintf("π%d", i)})
	}
	types.RegisterValue(&testGreeter{})

	// R1 gets boot value of R0
	r1b0, err := r[1].TryDial(l[0].Addr(), "test")
	if err != nil {
		t.Fatalf("dial r1b0 (%s)", err)
	}

	// R1 gets an ptr to a new greeter residing on R0
	g0 := r1b0.Call("NewGreeter")[0].(circuit.X)

	// R1 gets boot value of R2
	r1b2, err := r[1].TryDial(l[2].Addr(), "test")
	if err != nil {
		t.Fatalf("dial r1b2 (%s)", err)
	}

	// R1 passes the greeter g0 to R2; R2 will invoke a method in g0
	g := r1b2.Call("UseGreeter", g0)[0].(string)
	if g != testGreeting {
		t.Errorf("exp %s, got %s", testGreeting, g)
	}

	// Kick GC to initiate relinquishPtr from R2 to R1
	runtime.GC()
	runtime.Gosched()
	time.Sleep(time.Second)
}

type testFunc struct{}

func (testFunc) Func(msg string, p *TestData) *string {
	s := msg + " world " + strconv.Itoa(p.M["a"].(int))
	return &s
}

type TestData struct {
	X int
	M map[string]interface{}
	I interface{}
}

func TestGo(t *testing.T) {
	l1 := NewSandbox()
	New(l1)

	l2 := NewSandbox()
	r2 := New(l2)

	types.RegisterFunc(testFunc{})

	x := &TestData{
		X: 3,
		M: make(map[string]interface{}),
		I: (*int)(nil),
	}

	for i := 0; i < 10; i++ {
		x.M["a"] = i
		reply := r2.Go(l1.Addr(), testFunc{}, "hello", x)
		if len(reply) != 1 {
			t.Fatalf("missing reply")
		}
		s, ok := reply[0].(*string)
		if !ok {
			t.Fatalf("reply type")
		}
		if *s != "hello world "+strconv.Itoa(i) {
			t.Fatalf("reply value")
		}
	}
}
