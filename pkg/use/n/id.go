// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package n

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"strconv"

	"github.com/hoijui/circuit/pkg/use/errors"
)

var ErrParse = errors.NewError("parse")

// WorkerID represents the identity of a circuit worker process.
type WorkerID string

// String returns a cononical string representation of this worker ID.
func (r WorkerID) String() string {
	return string(r)
}

// ChooseWorkerID returns a random worker ID.
func ChooseWorkerID() WorkerID {
	return Int64WorkerID(rand.Int63())
}

func Int64WorkerID(src int64) WorkerID {
	return WorkerID(fmt.Sprintf("Q%016x", src))
}

func UInt64WorkerID(src uint64) WorkerID {
	return WorkerID(fmt.Sprintf("Q%016x", src))
}

// ParseOrHashWorkerID tries to parse the string s as a canonical worker ID representation.
// If it fails, it treats s as an unconstrained string and hashes it to a worker ID value.
// In either case, it returns a WorkerID value.
func ParseOrHashWorkerID(s string) WorkerID {
	id, err := ParseWorkerID(s)
	if err != nil {
		return HashWorkerID(s)
	}
	return id
}

// ParseWorkerID parses the string s for a canonical representation of a worker
// ID and returns a corresponding WorkerID value.
func ParseWorkerID(s string) (WorkerID, error) {
	if len(s) != 17 || s[0] != 'Q' {
		return "", ErrParse
	}
	ui64, err := strconv.ParseUint(s[1:], 16, 64)
	if err != nil {
		return "", err
	}
	return UInt64WorkerID(ui64), nil
}

// HashWorkerID hashes the unconstrained string s into a worker ID value.
func HashWorkerID(s string) WorkerID {
	h := fnv.New64a()
	h.Write([]byte(s))
	return UInt64WorkerID(h.Sum64())
}
