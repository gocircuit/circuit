// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/gocircuit/circuit/kit/interruptible"
	"github.com/gocircuit/circuit/use/circuit"
)

type Proc interface {
	Scrub()
	Wait() (Stat, error)
	Signal(sig string) error
	GetEnv() []string
	GetCmd() Cmd
	IsDone() bool
	Peek() Stat
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
	X() circuit.X
}

type proc struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	wait   <-chan error
	abr    <-chan struct{}
	cmd    struct {
		sync.Mutex
		cmd  exec.Cmd
		scrb bool
		abr  chan<- struct{}
		wait chan<- error
		exit error // exit set by waiter
	}
}

func MakeProc(cmd Cmd) Proc {
	p := &proc{}
	// std*
	p.cmd.cmd.Stdin, p.stdin = interruptible.BufferPipe(32e3)
	p.stdout, p.cmd.cmd.Stdout = interruptible.BufferPipe(32e3)
	p.stderr, p.cmd.cmd.Stderr = interruptible.BufferPipe(32e3)
	// exit
	ch, abr := make(chan error, 1), make(chan struct{})
	p.cmd.wait, p.wait = ch, ch
	p.abr, p.cmd.abr = abr, abr
	// cmd
	p.cmd.cmd.Env = cmd.Env
	p.cmd.cmd.Dir = cmd.Dir
	bin := strings.TrimSpace(cmd.Path)
	p.cmd.cmd.Path = bin
	p.cmd.cmd.Args = append([]string{bin}, cmd.Args...)
	p.cmd.scrb = cmd.Scrub
	// exec
	if err := p.cmd.cmd.Start(); err != nil {
		p.cmd.wait <- fmt.Errorf("exec error: %s", err.Error())
		close(p.cmd.wait)
		return p
	}
	go func() {
		p.cmd.wait <- p.cmd.cmd.Wait()
		close(p.cmd.wait)
		p.cmd.cmd.Stdout.(io.Closer).Close()
		p.cmd.cmd.Stderr.(io.Closer).Close()
	}()
	return p
}

func (p *proc) Stdin() io.WriteCloser {
	return p.stdin
}

func (p *proc) Stdout() io.ReadCloser {
	return p.stdout
}

func (p *proc) Stderr() io.ReadCloser {
	return p.stderr
}

func (p *proc) X() circuit.X {
	return circuit.Ref(XProc{p})
}

func (p *proc) Scrub() {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	if p.cmd.abr == nil {
		return
	}
	close(p.cmd.abr)
	p.cmd.abr = nil
}

func (p *proc) Wait() (Stat, error) {
	select {
	case exit, ok := <-p.wait:
		if !ok {
			return p.Peek(), nil
		}
		p.cmd.Lock()
		defer p.cmd.Unlock()
		p.cmd.exit = exit
		return p.peek(), nil
	case <-p.abr:
		return Stat{}, errors.New("aborted")
	}
}

func (p *proc) Signal(sig string) error {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	if p.cmd.cmd.Process == nil {
		return errors.New("no running process to signal")
	}
	if sig, ok := sigMap[strings.TrimSpace(strings.ToUpper(sig))]; ok {
		return p.cmd.cmd.Process.Signal(sig)
	}
	return errors.New("signal name not recognized")
}

func (p *proc) GetEnv() []string {
	return os.Environ()
}

func (p *proc) GetCmd() Cmd {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	return Cmd{
		Env:   p.cmd.cmd.Env,
		Path:  p.cmd.cmd.Path,
		Args:  p.cmd.cmd.Args[1:],
		Scrub: p.cmd.scrb,
	}
}

func (p *proc) IsDone() bool {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	if p.cmd.abr == nil {
		return true
	}
	switch p.phase() {
	case NotStarted, Exited, Signaled:
		return p.cmd.scrb
	}
	return false
}

func (p *proc) Peek() Stat {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	return p.peek()
}

func (p *proc) peek() Stat {
	return Stat{
		Cmd: Cmd{
			Env:   p.cmd.cmd.Env,
			Path:  p.cmd.cmd.Path,
			Args:  p.cmd.cmd.Args[1:],
			Scrub: p.cmd.scrb,
		},
		Exit:  p.cmd.exit,
		Phase: p.phase().String(),
	}
}

func (p *proc) phase() Phase {
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
