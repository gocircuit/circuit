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
	"path/filepath"
	"strings"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/interruptible"
)

type Select struct {
	*file.ErrorFile
	clauses struct {
		sync.Mutex
		waitfiles []string
	}
	back struct {
		interruptible.Mutex    // Both Select and Wait operations synchronize on this lock
		report  <-chan *result // waitOpenFile goroutines report open results to this channel
		abort   rh.Abandon     // closing intr interrupts all blocked waitOpenFile goroutines
		result  *result
	}
}

func MakeSelect() *Select {
	return &Select{ErrorFile: file.NewErrorFile()}
}

func (s *Select) WaitFiles() []string {
	s.clauses.Lock()
	defer s.clauses.Unlock()
	return s.clauses.waitfiles
}

// Select initiates a selection process.
func (s *Select) Select(waitfiles []string) (err error) {
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
	if s.clauses.waitfiles != nil {
		s.ErrorFile.Set("selection already in progress")
		return fmt.Errorf("select already running")
	}
	if waitfiles, err = verifyFiles(waitfiles); err != nil {
		s.ErrorFile.Set(err.Error())
		return err
	}
	s.clauses.waitfiles = waitfiles
	s.start(waitfiles)
	return nil
}

// verifyFiles runs some basic sanity checks on the wait file names
func verifyFiles(waitfiles []string) ([]string, error) {
	if len(waitfiles) == 0 {
		return nil, fmt.Errorf("no wait files")
	}
	for i, file := range waitfiles {
		file = path.Clean(file)
		_, filename := filepath.Split(file)
		if !strings.HasPrefix(filename, "wait") {
			return nil, fmt.Errorf("file %s is not a wait file", file)
		}
		waitfiles[i] = file
		// Make sure the files exist
		if _, err := os.Stat(file); err != nil {
			return nil, fmt.Errorf("cannot stat file %s: %s", file, err.Error())
		}
	}
	return waitfiles, nil
}

func (s *Select) start(waitfiles []string) {
	ch := make(chan *result, len(waitfiles))
	s.back.report = ch
	var intr rh.Intr
	intr, s.back.abort = rh.NewIntr()
	for i, file := range waitfiles {
		go waitOpenFile(i, file, intr, ch)
	}
}

// Wait returns when one of the wait files opens (with or without a POSIX open error).
// The POSIX file open error is not reflected in the final return value of Wait;
// The error returned by Wait is non-nil only in the case of an interruption or another system error.
func (s *Select) Wait(intr rh.Intr) (branch int, file string, err error) {
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
		return s.back.result.Branch, s.back.result.Name, nil
	}
	if s.back.report == nil {
		s.ErrorFile.Set("no selection in progress")
		return -1, "", rh.ErrClash
	}
	select {
	case s.back.result = <-s.back.report:
		s.back.report = nil // make sure no other results are consumed
		close(s.back.abort)  // kill all remaining waiters
		s.back.abort = nil
		return s.back.result.Branch, s.back.result.Name, nil
	case <-intr:
		s.ErrorFile.Set("wait interrupted")
		return -1, "", rh.ErrIntr
	}
}

// Clunk will abort the selection and terminate any waiting goroutines,
// as long as no Wait or Select operation is currently underway.
func (s *Select) Clunk() error {
	// obtain operation lock
	u := s.back.TryLock()
	if u == nil {
		return rh.ErrBusy
	}
	defer u.Unlock()
	//
	if s.back.abort == nil {
		return nil // either already done or not initialized yet
	}
	close(s.back.abort) // kill any outstanding file-open waiter processes
	s.back.abort = nil
	s.back.report = nil
	return nil
}

func (s *Select) TryWait() (branch int, file string, err error) {
	s.ErrorFile.Clear()
	u := s.back.TryLock()
	if u == nil {
		s.ErrorFile.Set("wait would block")
		return -1, "", rh.ErrBusy
	}
	defer u.Unlock()
	//
	if s.back.result != nil { // wait already unblocked
		return s.back.result.Branch, s.back.result.Name, nil
	}
	if s.back.report == nil {
		s.ErrorFile.Set("no selection in progress")
		return -1, "", rh.ErrClash
	}
	select {
	case s.back.result = <-s.back.report:
		s.back.report = nil // make sure no other results are consumed
		close(s.back.abort)  // kill all remaining waiters
		s.back.abort = nil
		return s.back.result.Branch, s.back.result.Name, nil
	default:
		s.ErrorFile.Set("no cases ready")
		return -1, "", rh.ErrBusy
	}
}
