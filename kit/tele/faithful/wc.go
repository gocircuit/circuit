// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package faithful provides a lossless chunked connection over a lossy one.
package faithful

import (
	"io"

	"github.com/gocircuit/circuit/kit/tele/chain"
	"github.com/gocircuit/circuit/kit/tele/trace"
)

type writeConn struct {
	frame trace.Frame
	sub   *chain.Conn
	bfr   *Buffer         // Buffer is self-synchronized
	sch   <-chan *control // Sync write channel; readLoop sends and closes, syncWriteLoop receives
}

// writeControl processes a single control message (SYNC or ACK).
// ok equals false if the connection is permanently broken and should be killed.
func writeControl(writer *chain.ConnWriter, msg encoder) {
	chunk, err := msg.Encode()
	if err != nil {
		panic(err)
	}
	writer.Write(chunk)
}

// writeUser processes a single set of return values from Buffer.Read.
func (wc *writeConn) writeUser(writer *chain.ConnWriter, payload interface{}, seqno SeqNo, err error) (continueWriteLoop bool) {
	if err == io.EOF {
		// We've reached the EOF of the user write sequence.
		// Go back to listen for more reads from the buffer, in case a sync rewinds the buffer cursor.
		return true
	}
	if err == io.ErrUnexpectedEOF {
		// An unexpected termination has been reached, indicated by killing the buffer; nothing to send any longer.
		return false
	}
	if err != nil {
		panic("u")
	}
	// Encode chunk
	chunk := payload.(*Chunk)
	chunk.seqno = seqno // Sequence numbers are assigned 0-based integers
	raw, err := chunk.Encode()
	if err != nil {
		panic(err)
	}
	writer.Write(raw)
	// If connection is closed (no more writes) and the buffer is empty, it is time to kill the connection.
	if wc.bfr.Drained() {
		return false
	}
	return true
}

func (wc *writeConn) loop() {
	bfrch := NewBufferReadChan(wc.bfr) // bfrchan returns a stream of chunks coming from buffer.Read
	defer func() {
		wc.bfr.Close()
		// Drain buffer Read until error
		for _ = range bfrch {
		}
		// Kill the underlying chain connection
		wc.sub.Kill()
	}()
	var writer *chain.ConnWriter
	for {
		// If no writer available, wait for one from the readLoop
		if writer == nil {
			ctrl, ok := <-wc.sch
			if !ok {
				return
			}
			if ctrl.Writer == nil {
				// Skip control messages that don't carry a new writer
				continue
			}
			writer = ctrl.Writer
			writeControl(writer, ctrl.Msg)
		}
		//
		select {
		case ctrl, ok := <-wc.sch:
			if !ok {
				return
			}
			if ctrl.Writer != nil {
				writer = ctrl.Writer
			}
			writeControl(writer, ctrl.Msg)

		case user, ok := <-bfrch:
			if !ok {
				return
			}
			if !wc.writeUser(writer, user.Payload, user.SeqNo, user.Err) {
				return
			}
		}
	}
}
