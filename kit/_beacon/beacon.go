// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package beacon

import (
	"net"
)

type 

// Discover…
func Discover(port int) <-chan string {
	ch := make(chan string)
	go func() {
		??
	}()
	go func() {
		// listen to broadcasts
		conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{ IP: net.IPv4bcast, Port: port })
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		// beam a request every 5 seconds until a response is heard
		for {
			??
			conn.Write([]byte("hola"))
		}
	}()
	return ch
}
