// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package faithful

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/gocircuit/circuit/kit/tele/chain"
)

// MsgKind determines the type of packet following
type MsgKind byte

const (
	CHUNK = MsgKind(iota)
	SYNC
	ACK
)

// Sequence number of a chunk sent over a HiFi chunk.Conn
type SeqNo int64

const MaxSeqNoVarintLen = binary.MaxVarintLen64

// Chunk is a message containing a chunk of user data.
type Chunk struct {
	seqno SeqNo
	chunk []byte
}

type encoder interface {
	Encode() ([]byte, error)
}

func (x *Chunk) String() string {
	return fmt.Sprintf("Chunk(SeqNo=%d, Len=%d)", x.seqno, len(x.chunk))
}

func (fh *Chunk) Encode() ([]byte, error) {
	var w bytes.Buffer
	w.WriteByte(byte(CHUNK))
	q := make([]byte, MaxSeqNoVarintLen)
	n := binary.PutVarint(q, int64(fh.seqno))
	w.Write(q[:n])
	// Writing the chunk's length is not necessary, since blobs are self-delimited.
	w.Write(fh.chunk)
	return w.Bytes(), nil
}

// Sync messages are sent by the receiver endpoint of the half-connection to request a retransmit.
type Sync struct {
	NAckd SeqNo
}

func (x *Sync) String() string {
	return fmt.Sprintf("Sync(NAckd=%d)", x.NAckd)
}

func (fh *Sync) Encode() ([]byte, error) {
	var w bytes.Buffer
	w.WriteByte(byte(SYNC))
	q := make([]byte, MaxSeqNoVarintLen)
	n := binary.PutVarint(q, int64(fh.NAckd))
	w.Write(q[:n])
	return w.Bytes(), nil
}

// Ack messages are sent by the receiver endpoint of the half-connection to announce what they have received.
type Ack struct {
	NAckd SeqNo
}

func (x *Ack) String() string {
	return fmt.Sprintf("Ack(NAckd=%d)", x.NAckd)
}

func (fh *Ack) Encode() ([]byte, error) {
	var w bytes.Buffer
	w.WriteByte(byte(ACK))
	q := make([]byte, MaxSeqNoVarintLen)
	n := binary.PutVarint(q, int64(fh.NAckd))
	w.Write(q[:n])
	return w.Bytes(), nil
}

// decodeMsg decodes either a *Chunk, *Sync or *Ack object from the byte array.
func decodeMsg(q []byte) (interface{}, error) {
	r := bytes.NewReader(q)
	// MsgKind
	msgkind, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	msgKind := MsgKind(msgkind)
	// Switch
	switch msgKind {
	case CHUNK:
		msg := &Chunk{}
		// SeqNo
		seqno, err := binary.ReadVarint(r)
		if err != nil {
			return nil, err
		}
		msg.seqno = SeqNo(seqno)
		// Chunk
		msg.chunk, err = ioutil.ReadAll(r)
		if err != nil {
			panic("u")
		}
		return msg, nil

	case SYNC:
		msg := &Sync{}
		nackd, err := binary.ReadVarint(r)
		if err != nil {
			return nil, err
		}
		msg.NAckd = SeqNo(nackd)
		return msg, nil

	case ACK:
		msg := &Ack{}
		nackd, err := binary.ReadVarint(r)
		if err != nil {
			return nil, err
		}
		msg.NAckd = SeqNo(nackd)
		return msg, nil
	}
	return nil, chain.ErrMisbehave
}
