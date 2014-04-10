// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"container/list"
)

type loop struct {
	cap      int
	stat     Stat
	queue    list.List
	waitsend list.List // chan<- struct{}
	waitrecv list.List // chan<- struct{}

	send  <-chan interface{}
	recv  <-chan chan<- interface{}
	sgone chan<- struct{}
	rgone chan<- struct{}
	req   <-chan interface{}
	resp  chan<- interface{}
}

func (l *loop) flushWaitSend() {
	defer l.waitsend.Init()
	for e := l.waitsend.Front(); e != nil; e = e.Next() {
		close(e.Value.(chan<- struct{}))
	}
}

func (l *loop) flushWaitRecv() {
	defer l.waitrecv.Init()
	for e := l.waitrecv.Front(); e != nil; e = e.Next() {
		close(e.Value.(chan<- struct{}))
	}
}

func (l *loop) main() {
	defer func() {
		close(l.resp)
		close(l.sgone)
		close(l.rgone)
		l.flushWaitSend()
		l.flushWaitRecv()
	}()
	for {
		var n = l.queue.Len()
		var err error
		switch {
		case n == 0:
			l.flushWaitSend()
			err = l.selectSend()
		case n >= l.cap:
			l.flushWaitRecv()
			err = l.selectRecv()
		default:
			l.flushWaitSend()
			l.flushWaitRecv()
			err = l.selectSendOrRecv()
		}
		if err != nil {
			return
		}
	}
}

func (l *loop) selectSend() error {
	select {
	case v := <-l.send:
		l.store(v)
		return nil
	case req := <-l.req:
		return l.reply(req)
	}
}

func (l *loop) selectRecv() error {
	select {
	case t := <-l.recv:
		l.transmit(t)
		return nil
	case req := <-l.req:
		return l.reply(req)
	}
}

func (l *loop) store(v interface{}) {
	l.queue.PushBack(v)
	l.stat.NSend++
}

func (l *loop) transmit(t chan<- interface{}) {
	f := l.queue.Front()
	l.queue.Remove(f)
	t <- f.Value // guaranteed instantaneous in Recv
	l.stat.NRecv++
}

func (l *loop) selectSendOrRecv() error {
	select {
	case v := <-l.send:
		l.store(v)
		return nil
	case t := <-l.recv:
		l.transmit(t)
		return nil
	case req := <-l.req:
		return l.reply(req)
	}
}

type reqSetCap struct {
	N int
}

type reqGetCap struct{}

type reqGetStat struct{}

type reqWaitSend struct {
	Unblock chan<- struct{}
}

type reqWaitRecv struct {
	Unblock chan<- struct{}
}

func (l *loop) reply(req interface{}) error {
	if req == nil {
		return ErrGone
	}
	switch t := req.(type) {
	case reqSetCap:
		l.cap = t.N
	case reqGetCap:
		l.resp <- l.cap
	case reqGetStat:
		var q Stat = l.stat // copy
		l.resp <- &q
	case reqWaitSend:
		l.waitsend.PushBack(t.Unblock)
	case reqWaitRecv:
		l.waitrecv.PushBack(t.Unblock)
	default:
		panic(0)
	}
	return nil
}
