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

type readConn struct {
	frame trace.Frame
	sub   *chain.Conn
	bfr   *Buffer            // TODO: Distinguish buffer sides // Seek, Remove, Abort
	rch   chan<- interface{} // []byte or error
	sch   chan<- *control    // Sync write channel; readLoop sends and closes, syncWriteLoop receives
	ach   <-chan struct{}    // Abort channel
	nread SeqNo              // Number of received chunks
	nackd SeqNo              // Number of sent and acknowledged chunks
}

// readConn

func (rc *readConn) syncWrite(writer *chain.ConnWriter, nackd SeqNo) {
	rc.sch <- &control{Writer: writer, Msg: &Sync{NAckd: nackd}}
}

func (rc *readConn) ackWrite(nackd SeqNo) {
	rc.sch <- &control{Writer: nil, Msg: &Ack{NAckd: nackd}}
}

func (rc *readConn) loop() {
	for {
		chunk, err := rc.read()
		if err != nil {
			// rc.frame.Printf("terminating: (%s)", err)
			rc.rch <- err // rch must have buffer cap for this final error in case no one is Reading
			break
		}
		rc.rch <- chunk
	}
	// Permanent connection end
	close(rc.rch)
	close(rc.sch)
	rc.bfr.Abort()
}

// read blocks until the next CHUNK is received. Meanwhile it processes incoming control messages like ACK and SYNC.
// Non-nil errors returned by read indicate irrecoverable physical errors on the underlying connection.
func (rc *readConn) read() ([]byte, error) {
	for {
		// Check for abort signal
		select {
		case <-rc.ach:
			return nil, io.ErrUnexpectedEOF
		default:
		}
		// The read cannot block on anything else other than reading on the underlying connection.
		// Note that sub.Kill will unblock sub.Read. The former will be called by writeLoop when
		// it receives the abortion signal.
		chunk, err := rc.sub.Read()

		// Stitching or permanent error
		if err != nil {
			if chunk != nil {
				panic("eh")
			}
			stitchConnWriter := chain.IsStitch(err)
			// Connection termination.
			// Any non-ErrStitch error implies connection termination.
			if stitchConnWriter == nil {
				// rc.frame.Println("read carrier error:", err.Error())
				return nil, err
			}
			// Stitching
			rc.syncWrite(stitchConnWriter, rc.nread)
			continue
		}

		// Payload received
		msg, err := decodeMsg(chunk)
		if err != nil {
			rc.frame.Println("read/decode error:", err.Error())
			// Misbehaved opponent is a connection termination.
			return nil, err
		}
		switch t := msg.(type) {
		case *Sync:
			if err = rc.readSync(t); err != nil {
				rc.frame.Println("read SYNC error:", err.Error())
				// Permanent connection termination.
				return nil, err
			}
			// Retry read

		case *Ack:
			if err = rc.readAck(t); err != nil {
				rc.frame.Println("read ACK error:", err.Error())
				// Permanent connection termination.
				return nil, err
			}
			// Retry read

		case *Chunk:
			if chunk := rc.readChunk(t); chunk != nil {
				return chunk, nil
			}
			// Redundant chunk was dropped. Retry read

		default:
			panic("eh")
		}
	}
}

// readChunk returns a non-nil chunk, if successful.
// Otherwise nil is returned to indicate that the packet was discarded.
func (rc *readConn) readChunk(chunkMsg *Chunk) []byte {
	// Is this a chunk that was already received?
	// Drop already-received duplicates.
	if chunkMsg.seqno < rc.nread {
		return nil
	}
	if chunkMsg.seqno == rc.nread {
		rc.nread++
		if rc.nread%AckFrequency == 0 {
			rc.ackWrite(rc.nread)
		}
		if chunkMsg == nil {
			panic("eh")
		}
		return chunkMsg.chunk
	}
	// Otherwise, a future packet implies lost packets. Request a sync. Drop the future packet.
	rc.syncWrite(nil, rc.nread)
	return nil
}

func (rc *readConn) readSync(syncMsg *Sync) error {
	//c.frame.Println("SYNC", syncMsg.NAckd)
	nackd := syncMsg.NAckd
	if nackd < rc.nackd || nackd > rc.bfr.NWritten() {
		return chain.ErrMisbehave
	}
	rc.nackd = nackd
	// Seek before Remove, so that new chunks don't race into the network
	// redundantly as a result of Remove, before the old ones have been resent.
	rc.bfr.Seek(nackd)
	rc.bfr.Remove(nackd)
	return nil
}

func (rc *readConn) readAck(ackMsg *Ack) error {
	//c.frame.Println("ACK", ackMsg.NAckd)
	nackd := ackMsg.NAckd
	if nackd < rc.nackd || nackd > rc.bfr.NWritten() {
		return chain.ErrMisbehave
	}
	rc.nackd = nackd
	rc.bfr.Remove(nackd)
	return nil
}
