// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package discover

import (
	"net"
	"testing"
)

func TestDiscovering(t *testing.T) {
	ch := make(chan int)
	maddr := &net.UDPAddr{IP: net.IP{228, 8, 8, 8}, Port: 8822}
	_, ch1 := New(maddr, []byte("d1"))
	_, ch2 := New(maddr, []byte("d2"))
	go func() {
		<-ch1
		ch <- 1
	}()
	go func() {
		<-ch2
		ch <- 1
	}()
	<-ch
	<-ch
}
