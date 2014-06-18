// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"fmt"
	"hash/fnv"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

// ReceiverID is a universal ID referring to the unique identity of a cross-reference receiver
type ReceiverID uint64

func (rid ReceiverID) String() string {
	return fmt.Sprintf("X%016x", uint64(rid))
}

func ChooseReceiverID() ReceiverID {
	return ReceiverID(rand.Int63())
}

var nonce ReceiverID // random nonce, unique to this process

func init() {
	rand.Seed(time.Now().UnixNano()) // main()'s seeding runs too late to capture this init
	nonce = ChooseReceiverID()
}

// ComputeReceiverID computes a unique across-processes ID for the receiver r
func ComputeReceiverID(r interface{}) ReceiverID {
	h := fnv.New64a()
	h.Write([]byte(nonce.String()))
	snapvalue(h, reflect.ValueOf(r))
	return ReceiverID(h.Sum64())
}

func snapvalue(w io.Writer, v reflect.Value) {
	var q string
	switch v.Kind() {

	case reflect.Bool: // Not addressable
		q = strconv.FormatBool(v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		q = strconv.FormatInt(v.Int(), 36)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		q = strconv.FormatUint(v.Uint(), 36)

	case reflect.Uintptr:
		panic("n/s")

	case reflect.Float32, reflect.Float64:
		q = strconv.FormatFloat(v.Float(), 'g', 65, 64)

	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		q = "(" + strconv.FormatFloat(real(c), 'g', 65, 64) + ", " + strconv.FormatFloat(imag(c), 'g', 65, 64) + "i)"

	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer: // Addressable
		q = strconv.FormatUint(uint64(v.Pointer()), 36) // uintptr

	case reflect.Interface:
		d := v.InterfaceData() // [2]uintptr
		q = "<" + strconv.FormatUint(uint64(d[0]), 36) + "," + strconv.FormatUint(uint64(d[1]), 36) + ">"

	case reflect.String:
		q = v.String()

	case reflect.Array:
		w.Write([]byte("{"))
		for i := 0; i < v.Len(); i++ {
			snapvalue(w, v.Index(i))
			w.Write([]byte(","))
		}
		w.Write([]byte("}"))
		return

	case reflect.Struct:
		w.Write([]byte("{"))
		for i := 0; i < v.NumField(); i++ {
			snapvalue(w, v.FieldByIndex([]int{i}))
			w.Write([]byte(","))
		}
		w.Write([]byte("}"))

	default:
		panic("u")
	}
	w.Write([]byte(q))
}
