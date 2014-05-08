// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	_ "github.com/gocircuit/circuit/kit/debug/ctrlc"
	_ "github.com/gocircuit/circuit/kit/debug/kill"

	"github.com/gocircuit/circuit/sys/lang"
	"github.com/gocircuit/circuit/sys/tele"

	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"

	"github.com/gocircuit/circuit/kit/kinfolk/locus"
)

// A one-off package side-effect initialization makes this process capable of talking to circuit workers.
func init() {
	rand.Seed(time.Now().UnixNano())
	tele.Init()
	t := n.NewTransport(n.ChooseWorkerID(), &net.TCPAddr{})
	fmt.Println(t.Addr().String())
	circuit.Bind(lang.New(t))
}

// Client is a live session with a circuit worker.
type Client struct {
	y locus.YLocus
}

func NewClient(workerURL string) *Client {
	c := &Client{}
	w, err := n.ParseAddr(workerURL)
	if err != nil {
		log.Fatalf("circuit address does not parse (%s)", err)
	}
	c.y = locus.YLocus{circuit.Dial(w, "locus")}
	return c
}

func (c *Client) Peers() []Terminal {
	peers := c.y.GetPeers()
	var r = make([]Terminal, len(peers))
	for i, p := range peers {
		r[i] = Terminal{y: p.Term, kin: p.Kin}
	}
	return r
}

// func (c *Client) newTerminal(xterm circuit.X, xkin ??) Terminal {
// 	return Terminal{
// 		y: ??,
// 		k: c.k
// 	}
// }