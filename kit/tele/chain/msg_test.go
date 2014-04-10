// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chain

import (
	"bytes"
	"reflect"
	"testing"
)

var (
	testDialMsg    = &msgDial{ID: 7, SeqNo: 1}
	testPayloadMsg = &msgPayload{Payload: []byte{0x7, 0x2, 0x3}}
)

func TestProtoDial(t *testing.T) {
	msg := testDialMsg
	var u bytes.Buffer
	if err := msg.Write(&u); err != nil {
		t.Fatalf("dial write (%s)", err)
	}
	dial, err := readMsgDial(&u)
	if err != nil {
		t.Fatalf("dial read (%s)", err)
	}
	if !reflect.DeepEqual(dial, msg) {
		t.Fatalf("expected %#v, got %#v", msg, dial)
	}
}

func TestProtoPayload(t *testing.T) {
	msg := testPayloadMsg
	var u bytes.Buffer
	if err := msg.Write(&u); err != nil {
		t.Fatalf("payload write (%s)", err)
	}
	payload, err := readMsgPayload(&u)
	if err != nil {
		t.Fatalf("payload read (%s)", err)
	}
	if !reflect.DeepEqual(payload, msg) {
		t.Fatalf("expected %#v, got %#v", msg, payload)
	}
}
