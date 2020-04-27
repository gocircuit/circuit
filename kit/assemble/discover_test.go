// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package assemble

import (
	//"fmt"
	"net"
	"testing"

	"github.com/hoijui/circuit/kit/xor"
)

// func TestAssembler(t *testing.T) {
// 	maddr := &net.UDPAddr{IP: net.IP{228, 8, 8, 8}, Port: 8822}
// 	a := NewAssembler(?, maddr)
// }

func TestDiscovering(t *testing.T) {
	ch := make(chan int)
	maddr := &net.UDPAddr{IP: net.IP{228, 8, 8, 8}, Port: 8822}
	scatter := NewScatter(maddr, xor.Key(0), []byte("d1"))
	gather := NewGatherLens(maddr, xor.Key(1), 2)
	go func() {
		scatter.Scatter()
	}()
	go func() {
		gather.Gather()
		ch <- 1
	}()
	<-ch
}
