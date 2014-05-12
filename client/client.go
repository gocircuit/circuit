// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"fmt"
	"errors"
	"math/rand"
	"net"
	"time"

	_ "github.com/gocircuit/circuit/kit/debug/ctrlc"
	_ "github.com/gocircuit/circuit/kit/debug/kill"

	"github.com/gocircuit/circuit/sys/lang"
	_ "github.com/gocircuit/circuit/sys/tele"

	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"

	"github.com/gocircuit/circuit/kit/anchor"
	"github.com/gocircuit/circuit/kit/kinfolk"
	"github.com/gocircuit/circuit/kit/kinfolk/locus"
)

// A one-off package side-effect initialization makes this process capable of talking to circuit workers.
func init() {
	rand.Seed(time.Now().UnixNano())
	t := n.NewTransport(n.ChooseWorkerID(), &net.TCPAddr{})
	fmt.Println(t.Addr().String())
	circuit.Bind(lang.New(t))
}

// Client is a live session with a circuit worker.
type Client struct {
	y locus.YLocus
}

func Dial(workerURL string) *Client {
	c := &Client{}
	w, err := n.ParseAddr(workerURL)
	if err != nil {
		panic("circuit address does not parse")
	}
	c.y = locus.YLocus{circuit.Dial(w, "locus")}
	return c
}

func (c *Client) Walk(walk []string) Anchor {
	if len(walk) == 0 {
		return c
	}
	p := c.y.GetPeers()[walk[0]]
	if p == nil {
		return nil
	}
	t := c.newTerminal(p.Term, p.Kin)
	return t.Walk(walk[1:])
}

func (c *Client) View() map[string]Anchor {
	var r = make(map[string]Anchor)
	for k, p := range c.y.GetPeers() {
		r[k] = c.newTerminal(p.Term, p.Kin)
	}
	return r
}

func (c *Client) newTerminal(xterm circuit.X, xkin kinfolk.KinXID) terminal {
	return terminal{
		y: anchor.YTerminal{xterm},
		k: xkin,
	}
}

func (c *Client) Worker() string {
	return ""
}

func (c *Client) MakeChan(n int) (Chan, error) {
	return nil, errors.New("cannot create elements outside of workers")
}

func (c *Client) MakeProc(cmd Cmd) (Proc, error) {
	return nil, errors.New("cannot create elements outside of workers")
}

func (c *Client) Get() interface{} {
	return nil
}
