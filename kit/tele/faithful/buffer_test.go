// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package faithful

import (
	"io"
	"testing"

	"github.com/gocircuit/circuit/kit/tele/trace"
)

type testSeqNo SeqNo

func TestBuffer(t *testing.T) {
	bfr := NewBuffer(trace.NewFrame("TestBuffer"), 2)
	bfr.Write(testSeqNo(0))
	bfr.Write(testSeqNo(1))
	bfr.Remove(1)
	bfr.Write(testSeqNo(2))
	bfr.Seek(SeqNo(1))
	// Read 1
	chunk, seqno, err := bfr.Read()
	if err != nil {
		t.Fatalf("read (%s) or bad seqno", err)
	}
	if chunk != testSeqNo(1) || seqno != SeqNo(1) {
		t.Fatalf("chunk=%d seqno=%d; expecting %d", chunk, seqno, 1)
	}
	// Read 2
	chunk, seqno, err = bfr.Read()
	if err != nil {
		t.Fatalf("read (%s)", err)
	}
	if chunk != testSeqNo(2) || seqno != SeqNo(2) {
		t.Fatalf("chunk=%d seqno=%d; expecting %d", chunk, seqno, 2)
	}
	bfr.Remove(2)
	bfr.Close()
	if _, _, err := bfr.Read(); err != io.EOF {
		t.Fatalf("u (%s)", err)
	}
}
