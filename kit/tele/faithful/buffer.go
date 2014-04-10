// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package faithful

import (
	"bytes"
	"container/list"
	"fmt"
	"io"
	"sync"

	"github.com/gocircuit/circuit/kit/tele/trace"
)

/*

	Buffer's invokation rules:
	==========================

	Abort, Remove, Seek must not compete.
	Remove or Seek must not be called after Abort.

	Write and Close must not compete.
	Neither Write nor Close can be called after Close.

	Read should be called greedily until it returns an error.

*/

// Buffer buffers up to a limit of consecutively numbered items.
// …
type Buffer struct {
	trace.Frame
	// Data structures
	x__      sync.Mutex // Linearize access to list and cursor
	list     list.List
	cursor   *list.Element
	seqno    SeqNo // Sequence number of the cursor element
	nremoved SeqNo // Number of removed writes
	nwritten SeqNo // Total number of Writes
	//
	wch chan struct{} // Write permissions
	// Read
	r__ sync.Mutex // Synchronize access to rch
	rch chan struct{}
	// Linearization locks
	ctlz__, wz__ sync.Mutex
}

// eof is an internal message type indicating end-of-transmission
type eof struct{}

// NewBuffer creates a new buffer with limit m.
func NewBuffer(frame trace.Frame, m int) *Buffer {
	b := &Buffer{
		wch: make(chan struct{}, m+1), // +1 for the final EOF, so it does not block
		// The capacity of rch is chosen so that writers to rch will never block.
		rch: make(chan struct{}, 4*m+2),
	}
	frame.Bind(b)
	b.Frame = frame
	for i := 0; i < m; i++ {
		b.wch <- struct{}{}
	}
	return b
}

func (b *Buffer) NWritten() SeqNo {
	b.x__.Lock()
	defer b.x__.Unlock()
	return b.nwritten
}

func (b *Buffer) IsEmpty() bool {
	b.x__.Lock()
	defer b.x__.Unlock()
	return b.list.Front() == nil
}

func (b *Buffer) String() string {
	b.x__.Lock()
	defer b.x__.Unlock()
	return b.StringNoLock()
}

func (b *Buffer) StringNoLock() string {
	var w bytes.Buffer
	for seqno, e := b.nremoved, b.list.Front(); e != nil; e, seqno = e.Next(), seqno+1 {
		if seqno > b.nremoved {
			w.WriteRune('•')
		}
		if e == b.cursor {
			fmt.Fprintf(&w, "(%d)", seqno)
		} else {
			fmt.Fprintf(&w, "%d", seqno)
		}
	}
	return string(w.Bytes())
}

// Write can block. It will return an io.ErrUnexpectedEOF if the buffer is closed.
func (b *Buffer) Write(v interface{}) error {
	b.wz__.Lock()
	defer b.wz__.Unlock()
	//
	// Obtain permission to write
	if _, ok := <-b.wch; !ok {
		close(b.rch)
		return io.ErrUnexpectedEOF
	}
	if readable, err := b.write(v); !readable || err != nil {
		return err
	}
	b.sendRead()
	return nil
}

// Close disallows future Writes and sends an EOF signal to the read-side.
// The read-side remains functional.
func (b *Buffer) Close() {
	b.wz__.Lock()
	defer b.wz__.Unlock()
	//
	if _, ok := <-b.wch; !ok {
		close(b.rch)
		return
	}
	if !b.writeEOF() { // Idempotent
		return
	}
	b.sendRead()
}

func (b *Buffer) write(v interface{}) (readable bool, err error) {
	b.x__.Lock()
	defer b.x__.Unlock()
	// Check that items are added in increasing order of sequence number and not after EOF
	if b.list.Len() > 0 {
		if _, isEOF := b.list.Back().Value.(eof); isEOF {
			return false, io.ErrClosedPipe
		}
	}
	b.list.PushBack(v)
	b.nwritten++
	if b.cursor == nil {
		b.cursor = b.list.Back()
		b.seqno = b.nwritten - 1
		return true, nil
	}
	return false, nil
}

func (b *Buffer) writeEOF() (readable bool) {
	b.x__.Lock()
	defer b.x__.Unlock()
	// Has EOF already been written? (EOF chunks are never removed from the buffer)
	back := b.list.Back()
	if back != nil {
		if _, isEOF := back.Value.(eof); isEOF {
			return false // No changes
		}
	}
	b.list.PushBack(eof{})
	b.nwritten++
	if b.cursor == nil {
		b.cursor = b.list.Back()
		b.seqno = b.nwritten - 1
		return true
	}
	return false
}

// sendRead sends a wake-up strobe to the read loop
func (b *Buffer) sendRead() {
	b.r__.Lock()
	defer b.r__.Unlock()
	b.rch <- struct{}{}
}

// Remove removes all items from the buffer whose sequence number is smaller than before.
func (b *Buffer) Remove(before SeqNo) {
	b.ctlz__.Lock()
	defer b.ctlz__.Unlock()
	//
	n, _ := b.remove(before)
	// Send write permissions
	for i := 0; i < n; i++ {
		b.wch <- struct{}{}
	}
}

// remove removes all chunks with sequence numbers smaller than before from the write memory.
// n equals the number of items removed from the buffer.
// changed is true if the current reading position was changed to a non-nil element;
// The cursor position can only change if it was non-nil to begin with.
func (b *Buffer) remove(before SeqNo) (n int, changed bool) {
	b.x__.Lock()
	defer b.x__.Unlock()
	for seqno, e := b.nremoved, b.list.Front(); e != nil && seqno < before; seqno++ {
		// Never delete the last EOF chunk, if present
		if _, isEOF := e.Value.(eof); isEOF {
			break
		}
		// Slide forward
		save := e.Next()
		// Slide cursor forward as well.
		if b.cursor == e {
			b.cursor = save
			b.seqno++
			if b.cursor != nil {
				changed = true
			}
		}
		b.list.Remove(e)
		n++          // Number of elements removed during this call
		b.nremoved++ // Number of elements removed ever
		e = save
	}
	return n, changed
}

// Seek moves the current position to the item in the buffer numbered seqno.
// Seek panics if the requested seqno is after the last message written.
func (b *Buffer) Seek(seqno SeqNo) {
	b.ctlz__.Lock()
	defer b.ctlz__.Unlock()
	//
	if !b.seek(seqno) {
		return
	}
	// If the buffer has been Closed, Seek should continue to work so that
	// Read can empty out the buffer.
	b.sendRead()
}

// seek panics if the requested seqno is after the last message written
// seek returns true if the resulting cursor position is readable.
func (b *Buffer) seek(seqno SeqNo) (readable bool) {
	b.x__.Lock()
	defer b.x__.Unlock()
	if seqno > b.nwritten {
		panic("seeking to unwritten message")
	}
	for iter, e := b.nremoved, b.list.Front(); e != nil; iter, e = iter+1, e.Next() {
		switch {
		case iter < seqno:
		case iter == seqno:
			b.cursor, b.seqno = e, iter
			return true
		default:
			panic("u")
		}
	}
	b.cursor, b.seqno = nil, 0
	return false
}

// Read returns the next item that is due in the read order.
// An io.EOF is returned if the next chunk is EOF, indicating that the stream of Writes has been closed, but the Buffer object is still operational.
// An io.ErrUnexpectedEOF is returned if the Buffer has been Kill'ed.
func (b *Buffer) Read() (interface{}, SeqNo, error) {
	for {
		chunk, seqno := b.read()
		if chunk != nil {
			if _, isEOF := chunk.(eof); isEOF {
				return nil, seqno, io.EOF
			}
			return chunk, seqno, nil
		}
		// If nothing read, wait for the next read strobe or kill signal.
		// Superflous read signals (when no new reads are available) are OK.
		if _, readReady := <-b.rch; !readReady {
			return nil, 0, io.ErrUnexpectedEOF
		}
	}
}

// read returns the next element in the buffer read order, or nil otherwise,
// and slides the cursor forward.
func (b *Buffer) read() (v interface{}, seqno SeqNo) {
	b.x__.Lock()
	defer b.x__.Unlock()
	if b.cursor == nil {
		return nil, 0
	}
	//b.Printf("BFR BEFORE %s", b.StringNoLock())
	v, seqno, b.cursor, b.seqno = b.cursor.Value, b.seqno, b.cursor.Next(), b.seqno+1
	//b.Printf("BFR AFTER  %s", b.StringNoLock())
	return v, seqno
}

// Abort

func (b *Buffer) Abort() {
	b.ctlz__.Lock()
	defer b.ctlz__.Unlock()
	//
	close(b.wch)
}

// Drained returns true if no more writes are allowed on the connection and the buffer is empty or has only an EOF chunk in it.
func (b *Buffer) Drained() bool {
	b.x__.Lock()
	defer b.x__.Unlock()
	if b.list.Len() != 1 {
		return false
	}
	if _, isEOF := b.list.Back().Value.(eof); isEOF {
		return true
	}
	return false
}
