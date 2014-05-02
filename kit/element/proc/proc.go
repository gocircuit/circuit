// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

type Proc struct {
	Stdin  interruptible.Writer
	Stdout interruptible.Reader
	Stderr interruptible.Reader
	wait struct {
		interruptible.Mutex
		wait <-chan error
	}
	cmd struct {
		sync.Mutex
		cmd exec.Cmd
		wait chan<- error
		exit error // exit set by waiter
	}
	*file.ErrorFile
}

func MakeProc() *Proc {
	p := &Proc{ErrorFile: file.NewErrorFile()}

	// stdin
	inr, inw := interruptible.Pipe()
	in, err := p.cmd.cmd.StdinPipe()
	if err != nil {
		panic(0)
	}
	p.Stdin = inw
	go func() {
		io.Copy(in, inr)
		in.Close()
	}()
	// stdout
	outr, outw := interruptible.Pipe()
	out, err := p.cmd.cmd.StdoutPipe()
	if err != nil {
		panic(0)
	}
	p.Stdout = outr
	go func() {
		io.Copy(outw, out)
		outw.Close()
	}()
	// stderr
	er, ew := interruptible.Pipe()
	e, err := p.cmd.cmd.StderrPipe()
	if err != nil {
		panic(0)
	}
	p.Stderr = er
	go func() {
		io.Copy(ew, e)
		ew.Close()
	}()


	ch := make(chan error, 1)
	p.cmd.wait, p.wait.wait = ch, ch
	return p
}

func (p *Proc) ClunkIfNotBusy() error {
	p.ErrorFile.Clear()
	p.cmd.Lock()
	defer p.cmd.Unlock()
	//
	switch p.getState() {
	case Unknown, None, Exited:
		return nil
	}
	p.ErrorFile.Set("process is busy")
	return rh.ErrBusy
}

func (p *Proc) Start() error {
	p.ErrorFile.Clear() // clear error file
	p.cmd.Lock()
	defer p.cmd.Unlock()
	//
	if err := p.cmd.cmd.Start(); err != nil {
		p.ErrorFile.Setf("exec error: %s", err.Error())
		return err
	}
	go func() {
		p.cmd.wait <- p.cmd.cmd.Wait()
		close(p.cmd.wait)
		println("bye")
	}()
	return nil
}

func (p *Proc) Wait(intr rh.Intr) (Stat, error) {
	p.ErrorFile.Clear()
	u := p.wait.Lock(intr)
	if u == nil {
		p.ErrorFile.Set("wait interrupted")
		return Stat{}, rh.ErrIntr
	}
	defer u.Unlock()
	//
	select {
	case err, ok := <-p.wait.wait:
		if !ok {
			p.ErrorFile.Set("process already exited")
			return p.TryWait(), nil
		}
		p.cmd.Lock()
		defer p.cmd.Unlock()
		p.cmd.exit = err
		return p.stat(), nil
	case <-intr:
		p.ErrorFile.Set("wait interrupted")
		return Stat{}, rh.ErrIntr
	}
}

func (p *Proc) Signal(sig string) error {
	p.ErrorFile.Clear() // clear error file
	p.cmd.Lock()
	defer p.cmd.Unlock()
	//
	if p.cmd.cmd.Process == nil {
		p.ErrorFile.Set("no running process to signal")
		return rh.ErrClash
	}
	if sig, ok := sigMap[strings.TrimSpace(strings.ToUpper(sig))]; ok {
		return p.cmd.cmd.Process.Signal(sig)
	}
	p.ErrorFile.Set("signal name/number not recognized")
	return rh.ErrClash
}

func (p *Proc) GetEnv() []string {
	return os.Environ()
}

func (p *Proc) SetCmd(i *Cmd) {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	//
	p.cmd.cmd.Env = i.Env
	bin := strings.TrimSpace(i.Path)
	p.cmd.cmd.Path = bin
	p.cmd.cmd.Args = append([]string{bin}, i.Args...)
}

func (p *Proc) GetCmd() *Cmd {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	//
	return &Cmd{
		Env:  p.cmd.cmd.Env,
		Path: p.cmd.cmd.Path,
		Args: p.cmd.cmd.Args,
	}
}

func (p *Proc) TryWait() Stat {
	p.cmd.Lock()
	defer p.cmd.Unlock()
	return p.stat()
}

func (p *Proc) stat() Stat {
	return Stat{
		Cmd: Cmd{
			Env: p.cmd.cmd.Env,
			Path: p.cmd.cmd.Path,
			Args: p.cmd.cmd.Args,
		},
		Exit: p.cmd.exit,
		State: p.getState().String(),
	}
}

func (p *Proc) getState() RunState {
	if p.cmd.cmd.Process == nil {
		return None
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
	return Unknown
}
