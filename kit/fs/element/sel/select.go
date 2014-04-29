// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package sel

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/interruptible"
)

type Select struct {
	dir *SelectDir
	*file.ErrorFile
	clauses struct {
		sync.Mutex
		clauses []Clause
	}
	back struct {
		interruptible.Mutex    // Select and Wait synchronize on this lock
		unblock <-chan *unblock // file-open goroutines report results to this channel
		abort   rh.Abandon      // closing intr interrupts all blocked file-open goroutines
		result  *unblock
	}
}

func MakeSelect(dir *SelectDir) *Select {
	return &Select{
		dir:  dir,
		ErrorFile: file.NewErrorFile(),
	}
}

func (s *Select) GetClauses() []Clause {
	s.clauses.Lock()
	defer s.clauses.Unlock()
	return s.clauses.clauses
}

// Select initiates a selection process.
func (s *Select) Select(clauses []Clause) (err error) {
	s.ErrorFile.Clear()
	// obtain operation lock
	u := s.back.TryLock()
	if u == nil {
		s.ErrorFile.Set("another select or wait operation is blocking")
		return rh.ErrBusy
	}
	defer u.Unlock()
	//
	s.clauses.Lock()
	defer s.clauses.Unlock()
	if s.clauses.clauses != nil {
		s.ErrorFile.Set("selection already in progress")
		return fmt.Errorf("selection already in progress")
	}
	if clauses, err = verifyFiles(clauses); err != nil {
		s.ErrorFile.Set(err.Error())
		return err
	}
	s.clauses.clauses = clauses
	s.start(clauses)
	return nil
}

// verifyFiles runs some basic sanity checks on the wait file names
func verifyFiles(clauses []Clause) ([]Clause, error) {
	if len(clauses) == 0 {
		return nil, fmt.Errorf("no clauses")
	}
	for i, clause := range clauses {
		clause.Op = strings.ToLower(strings.TrimSpace(clause.Op))
		switch clause.Op {
		case "r", "w":
		default:
			return nil, fmt.Errorf("unknown operation")
		}
		clause.File = path.Clean(clause.File)
		clauses[i] = clause
		if _, err := os.Stat(clause.File); err != nil { // Make sure the files exist
			return nil, fmt.Errorf("cannot stat file %s: %s", clause.File, err.Error())
		}
	}
	return clauses, nil
}

func (s *Select) start(clauses []Clause) {
	ch := make(chan *unblock, len(clauses))
	s.back.unblock = ch
	var intr rh.Intr
	intr, s.back.abort = rh.NewIntr()
	for i_, clause_ := range clauses {
		i, clause := i_, clause_
		switch clause.Op {
		case "r":
			go func() {
				u := &unblock{Clause: i}
				defer func() { // catch panics from disappearing clause files
					if r := recover(); r != nil {
						u.Error = rh.ErrNotExist
						ch <- u
					}
				}()
				u.File, u.Error = OpenFileReader(clause.File, intr)
				ch <- u
			}()
		case "w":
			go func() {
				u := &unblock{Clause: i}
				defer func() { // catch panics from disappearing clause files
					if r := recover(); r != nil {
						u.Error = rh.ErrNotExist
						ch <- u
					}
				}()
				u.File, u.Error = OpenFileWriter(clause.File, intr)
				ch <- u
			}()
		}
	}
}

type unblock struct {
	Clause int
	File   interface{} // *FileReader or *FileWriter
	Error  error
}

func (u *unblock) CommitName() string {
	if u.File == nil {
		return ""
	}
	switch u.File.(type) {
	case *FileReader:
		return "read." + strconv.Itoa(u.Clause)
	case *FileWriter:
		return "write." + strconv.Itoa(u.Clause)
	}
	panic(0)
}

func (u *unblock) Return() (clause int, commit string, err error) {
	return u.Clause, u.CommitName(), u.Error
}

// Wait returns when one of the wait files opens (with or without a POSIX open error).
// The error returned by Wait is non-nil only in the case of an interruption or another system error.
func (s *Select) Wait(intr rh.Intr) (clause int, commit string, err error) {
	s.ErrorFile.Clear()
	// obtain operation lock
	u := s.back.Lock(intr)
	if u == nil {
		s.ErrorFile.Set("wait interrupted")
		return -1, "", rh.ErrIntr
	}
	defer u.Unlock()
	//
	if s.back.result != nil { // wait already unblocked
		return s.back.result.Return()
	}
	if s.back.unblock == nil {
		s.ErrorFile.Set("no selection in progress")
		return -1, "", rh.ErrClash
	}
	//
	select {
	case s.back.result = <-s.back.unblock:
		s.abort() // stop all other waiters
		s.plant()
		return s.back.result.Return()
	case <-intr:
		s.ErrorFile.Set("wait interrupted")
		return -1, "", rh.ErrIntr
	}
}

func (s *Select) plant() {
	r := s.back.result
	if r.Error != nil {
		return // no file to plant if unblock was an error
	}
	switch t := r.File.(type) {
	case *FileReader:
		s.dir.dir.AddChild(r.CommitName(), file.NewFileFID(NewDelayedReadFile(t)))
	case *FileWriter:
		s.dir.dir.AddChild(r.CommitName(), file.NewFileFID(NewDelayedWriteFile(t)))
	}
	panic(0)
}

func (s *Select) abort() error {
	if s.back.abort == nil {
		return rh.ErrGone
	}
	s.back.unblock = nil // make sure no other unblock results are consumed
	close(s.back.abort)  // kill all remaining waiters
	runtime.GC()         // rush the collection of *FileReaders and *FileWriters
	s.back.abort = nil
	return nil
}

// Abort will abort the selection and terminate any waiting goroutines,
// as long as no Wait or Select operation is currently underway.
func (s *Select) Abort() (err error) {
	// obtain operation lock
	u := s.back.TryLock()
	if u == nil {
		return rh.ErrBusy
	}
	defer u.Unlock()
	//
	return s.abort()
}

// Scrub is like abort, but it will fail with rh.ErrBusy if Wait or Select is pending.
func (s *Select) Scrub() (err error) {
	// obtain operation lock
	u := s.back.TryLock()
	if u == nil {
		return rh.ErrBusy
	}
	defer u.Unlock()
	//
	err = s.abort()
	switch err {
	case rh.ErrGone, nil:
		return nil
	default:
		return err
	}
}

func (s *Select) TryWait() (clause int, commit string, err error) {
	s.ErrorFile.Clear()
	u := s.back.TryLock()
	if u == nil {
		s.ErrorFile.Set("wait would block")
		return -1, "", rh.ErrBusy
	}
	defer u.Unlock()
	//
	if s.back.result != nil { // wait already unblocked
		return s.back.result.Return()
	}
	if s.back.unblock == nil {
		s.ErrorFile.Set("no selection in progress")
		return -1, "", rh.ErrClash
	}
	//
	select {
	case s.back.result = <-s.back.unblock:
		s.abort() // stop all other waiters
		s.plant()
		return s.back.result.Return()
	default:
		s.ErrorFile.Set("no clauses ready")
		return -1, "", rh.ErrBusy
	}
}
