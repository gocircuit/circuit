// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

/*
	Package client provides access to the circuit programming environment to user programs.
*/
package client

import (
	"errors"
	"math/rand"
	"net"
	"time"
	"sync"

	_ "github.com/gocircuit/circuit/kit/debug/kill"

	"github.com/gocircuit/circuit/sys/lang"
	_ "github.com/gocircuit/circuit/sys/tele"

	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/kinfolk"
	"github.com/gocircuit/circuit/kinfolk/locus"
)

var _once sync.Once
func _init(key []byte) {
	rand.Seed(time.Now().UnixNano())
	t := n.NewTransport(n.ChooseWorkerID(), &net.TCPAddr{}, key)
	//fmt.Println(t.Addr().String())
	circuit.Bind(lang.New(t))
}

// Client is a live session with a circuit server.
type Client struct {
	y locus.YLocus
}

// DialAddr establishes a connection to a circuit server specified by a circuit address.
// Circuit addresses are printed to standard output when a server is started with the 
// "circuit start …" command.
//
// If key is non-nil it is used as a private key and all communications are
// secured by HMAC authentication and RC4 symmetric encryption;
// otherwise transmissions are in plaintext.
//
// Errors in communication, such as a missing server, or invalid URL format
// are reported through panics.
func Dial(addr string, key []byte) *Client {
	_once.Do(func() {
		_init(key)
	})
	c := &Client{}
	w, err := n.ParseAddr(addr)
	if err != nil {
		panic("circuit address does not parse")
	}
	c.y = locus.YLocus{circuit.Dial(w, "locus")}
	return c
}

// DialDiscover…
// func DialDiscover(addr string, key []byte) *Client {
// 	_once.Do(func() {
// 		_init(key)
// 	})
// 	c := &Client{}
// 	?
// }

// Address returns the circuit address of the server that this client is connected to.
func (c *Client) Addr() string {
	return c.y.X.Addr().String()
}

// Walk traverses the global virtual anchor namespace and returns a handle dor the desired anchor.
// The first element of walk should be the ID of a live circuit server.
// An up to date list of available circuit servers in the cluster can be obtained by calling View.
// The remainder of the walk slice is up to the user.
// Errors in communication or missing servers are reported as panics.
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

// View returns a map of all currently-live circuit server anchors.
// Errors in communication are reported as panics.
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

// ServerID is an Anchor interface method, not applicable to the root-level anchor.
func (c *Client) ServerID() string {
	return ""
}

// MakeChan is an Anchor interface method, not applicable to the root-level anchor.
func (c *Client) MakeChan(n int) (Chan, error) {
	return nil, errors.New("cannot create elements outside of workers")
}

// MakeProc is an Anchor interface method, not applicable to the root-level anchor.
func (c *Client) MakeProc(cmd Cmd) (Proc, error) {
	return nil, errors.New("cannot create elements outside of workers")
}

// MakeOnJoin is an Anchor interface method, not applicable to the root-level anchor.
func (c *Client) MakeOnJoin() (Subscription, error) {
	return nil, errors.New("cannot create elements outside of workers")
}

// MakeOnLeave is an Anchor interface method, not applicable to the root-level anchor.
func (c *Client) MakeOnLeave() (Subscription, error) {
	return nil, errors.New("cannot create elements outside of workers")
}

// Get is an Anchor interface method, not applicable to the root-level anchor.
func (c *Client) Get() interface{} {
	return nil
}

// Scrub is an Anchor interface method, not applicable to the root-level anchor.
func (c *Client) Scrub() {}