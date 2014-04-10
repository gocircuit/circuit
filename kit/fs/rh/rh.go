// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package rh provides a the interface for the Resource Hierarchy file system
package rh

// Originally inspired from:
//	http://plan9.bell-labs.com/sys/man/5/INDEX.html

// Server is the low-level interface that a X9P file server must implement
type Server interface {

	// SignIn opens a new session with the server for user, rooted at dir
	SignIn(user, dir string) (Session, error)

	// String returns a short name, identifying this server
	String() string
}

type Session interface {

	// Walk returns a new unopened FID, corresponding to the given relative path.
	// In particular, Walk(nil) returns an unopened FID for the root.
	Walk(name []string) (FID, error)

	// ?
	SignOut()

	// String returns a short name, identifying this session
	String() string
}

type Identifier interface {

	// String returns a short readable name, identifying this FID. The name need not conform to any syntactic rules.
	String() string

	// Q returns the QID of this FID.
	Q() Q
}

type Walker interface {
	// Walk returns a new unopened FID, corresponding to the given relative path.
	// In particular, Walk(nil) returns an unopened FID for the current file.
	Walk(wname []string) (FID, error)
}

type Stater interface {

	// Stat inquires about the file identified by this FID.
	Stat() (*Dir, error)

	// Wstat can change some of the file's meta information.
	Wstat(*Wdir) error
}

// Prompt is the value of an interrupt notification.
// Currently unused, it is meant for things like OS signal info, etc.
type Prompt interface{}

// Intr is a channel that receives interrupt notifications
type Intr <-chan Prompt

// Abandon is a channel that sends interrupt notifications
type Abandon chan<- Prompt

func NewIntr() (Intr, Abandon) {
	ch := make(chan Prompt, 1)
	return Intr(ch), Abandon(ch)
}

//
type Conn interface {

	// Open asks the file server to check permissions and prepare this FID for I/O with subsequent read and write invocations.
	Open(flag Flag, intr Intr) error

	//
	Create(name string, flag Flag, mode Mode, perm Perm) (FID, error)

	// Clunk releases this FID immediately.
	// The FID should also be released automatically if the FID object is garbage-collected.
	Clunk() error

	// If count equals zero, all data should be returned.
	Read(offset int64, count int, intr Intr) (Chunk, error)

	//
	Write(offset int64, data Chunk, intr Intr) (int, error)
}

type Traverser interface {

	// Remove asks the file server both to remove this file and to clunk this FID, even if the remove fails.
	// If a file has been opened as multiple fids, possibly on different
	// connections, and one fid is used to remove the file, whether the other
	// fids continue to provide access to the file is implementation–defined.
	Remove() error

	// Move moves this FID to dir with the given name.
	Move(dir FID, name string) error
}

type FID interface {
	Identifier
	Walker
	Stater
	Conn
	Traverser
}

const (
	MaxWalk = 16 // Suggested maximum depth of walk sequences
)

// Mixing AnchoredFID into your type takes care of rejecting meta- or location-modifying requests.
type AnchoredFID struct{}

func (AnchoredFID) Create(string, Flag, Mode, Perm) (FID, error) {
	return nil, ErrPerm
}

func (AnchoredFID) Wstat(*Wdir) error {
	return ErrPerm
}

func (AnchoredFID) Remove() error {
	return ErrPerm
}

func (AnchoredFID) Move(FID, string) error {
	return ErrPerm
}

// Mixing ReadOnlyFID into your type takes care of rejecting write-side calls.
type ReadOnlyFID struct {
	AnchoredFID
	DontWriteFID
}

type DontWriteFID struct{}

func (DontWriteFID) Write(int64, Chunk, Intr) (int, error) {
	return 0, ErrPerm
}

// Mixing WriteOnlyFID into your type takes care of rejecting read-side calls.
type WriteOnlyFID struct {
	AnchoredFID
	DontReadFID
}

type DontReadFID struct{}

func (DontReadFID) Read(int64, int, Intr) (Chunk, error) {
	return nil, ErrPerm
}

//
type ZeroFID struct {
	AnchoredFID
	DontReadFID
	DontWriteFID
}
