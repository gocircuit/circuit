// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/gocircuit/circuit/kit/interruptible"
)

type Proc struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	wait <-chan error
	cmd struct {
		sync.Mutex
		cmd exec.Cmd
		wait chan<- error
		exit error // exit set by waiter
	}
}

func MakeProc(cmd *Cmd) *Proc {
	var err error
	p := &Proc{}
	// std*
	if p.Stdin, err = p.cmd.cmd.StdinPipe(); err != nil {
		panic(0)
	}
	if p.Stdout, err = p.cmd.cmd.StdoutPipe(); err != nil {
		panic(0)
	}
	if p.Stderr, err = p.cmd.cmd.StderrPipe(); err != nil {
		panic(0)
	}
	// exit
	ch := make(chan error, 1)
	p.cmd.wait, p.wait = ch, ch
	// cmd
	p.cmd.cmd.Env = cmd.Env
	bin := strings.TrimSpace(cmd.Path)
	p.cmd.cmd.Path = bin
	p.cmd.cmd.Args = append([]string{bin}, cmd.Args...)
	// exec
	if err := p.cmd.cmd.Start(); err != nil {
		p.cmd.wait <- fmt.Errorf("exec error: %s", err.Error())
		close(p.cmd.wait)
		return p
	}
	go func() {
		p.cmd.wait <- p.cmd.cmd.Wait()
		close(p.cmd.wait)
	}()
	return p
}

func (p *Proc) Wait(intr interruptible.Intr) (Stat, error) {
	select {
	case err, ok := <-p.wait:
		if !ok {
			return p.Peek(), nil
		}
		p.cmd.Lock()
		defer p.cmd.Unlock()
		p.cmd.exit = err
		return p.peek(), nil
	case <-intr:
		return Stat{}, interruptible.ErrIntr
	}
}

func (p *Proc) Signal(sig string) error {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	if p.cmd.cmd.Process == nil {
		return fmt.Errorf("no running process to signal")
	}
	if sig, ok := sigMap[strings.TrimSpace(strings.ToUpper(sig))]; ok {
		return p.cmd.cmd.Process.Signal(sig)
	}
	return fmt.Errorf("signal name not recognized")
}

func (p *Proc) GetEnv() []string {
	return os.Environ()
}

func (p *Proc) GetCmd() *Cmd {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	return &Cmd{
		Env:  p.cmd.cmd.Env,
		Path: p.cmd.cmd.Path,
		Args: p.cmd.cmd.Args,
	}
}

func (p *Proc) IsDone() bool {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	switch p.phase() {
	case NotStarted, Exited, Signaled:
		return true
	}
	return false
}

func (p *Proc) Peek() Stat {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	return p.peek()
}

func (p *Proc) peek() Stat {
	return Stat{
		Cmd: Cmd{
			Env: p.cmd.cmd.Env,
			Path: p.cmd.cmd.Path,
			Args: p.cmd.cmd.Args,
		},
		Exit: p.cmd.exit,
		Phase: p.phase().String(),
	}
}

func (p *Proc) phase() Phase {
	if p.cmd.cmd.Process == nil {
		return NotStarted // didn't start due to error
	}
	ps := p.cmd.cmd.ProcessState
	if ps == nil {
		return Running
	}
	ws := ps.Sys().(syscall.WaitStatus)
	switch {
	case ps.Exited():
		return Exited
	case ws.Stopped():
		return Stopped
	case ws.Signaled():
		return Signaled
	case ws.Continued():
		return Continued
	}
	panic(0)
}
