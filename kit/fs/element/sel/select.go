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
	sel struct {
		sync.Mutex
		waitfiles []string
	}
	wait struct {
		l1 interruptible.Mutex
		l2 sync.Mutex
		ch <-chan result
	}
	*file.ErrorFile
}

func MakeSelect() *Select {
	return &Select{ErrorFile: file.NewErrorFile()}
}

func (s *Select) ClunkIfNotBusy() error {
	s.wait.l2.Lock()
	ch := s.wait.ch
	s.wait.l2.Unlock()
	if ch == nil {
		return nil
	}
	return rh.ErrBusy
}

func (s *Select) Wait(intr rh.Intr) (file string, err error) {
	s.ErrorFile.Clear()
	u := s.wait.l1.Lock(intr)
	if u == nil {
		s.ErrorFile.Set("wait interrupted")
		return "", rh.ErrIntr
	}
	defer u.Unlock()
	//
	s.wait.l2.Lock()
	ch := s.wait.ch
	s.wait.l2.Unlock()
	//
	if ch == nil {
		s.ErrorFile.Set("no selection in progress")
		return "", rh.ErrClash
	}
	select {
	case r := <-ch:
		// reset the select for reuse
		s.sel.Lock()
		defer s.sel.Unlock()
		//
		s.wait.l2.Lock()
		defer s.wait.l2.Unlock()
		//
		s.sel.waitfiles = nil
		s.wait.ch = nil
		//
		return r.File, r.Error
	case <-intr:
		s.ErrorFile.Set("wait interrupted")
		return "", rh.ErrIntr
	}
}

func (s *Select) WaitFiles() []string {
	s.sel.Lock()
	defer s.sel.Unlock()
	return s.sel.waitfiles
}

func (s *Select) Select(waitfiles []string) (err error) {
	s.ErrorFile.Set("") // clear error file
	s.sel.Lock()
	defer s.sel.Unlock()
	if s.sel.waitfiles != nil {
		s.ErrorFile.Set("selection already in progress")
		return fmt.Errorf("select already running")
	}
	if waitfiles, err = verifyFiles(waitfiles); err != nil {
		s.ErrorFile.Set(err.Error())
		return err
	}
	s.sel.waitfiles = waitfiles
	s.start(waitfiles)
	return nil
}

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
	ch := make(chan result, len(waitfiles))
	retrn := &returnResult{ch: ch}
	for _, file := range waitfiles {
		go waitOpen(file, retrn)
	}
	s.wait.l2.Lock()
	defer s.wait.l2.Unlock()
	s.wait.ch = ch
}

type result struct {
	File  string
	Error error
}

type returnResult struct {
	sync.Mutex
	ch chan<- result
}

func (rr *returnResult) Return(file string, err error) {
	rr.Lock()
	defer rr.Unlock()
	rr.ch <- result{file, err}
}

func waitOpen(file string, retrn *returnResult) {
	f, err := os.Open(file)
	if err != nil {
		retrn.Return(file, err)
		return
	}
	defer f.Close()
	retrn.Return(file, nil)
}
