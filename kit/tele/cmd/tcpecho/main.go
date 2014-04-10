// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"flag"
	"log"
	"net"
	"os"
)

var flagAddr = flag.String("addr", ":8787", "Address to listen to")

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", *flagAddr)
	if err != nil {
		log.Printf("accept (%s)", err)
		os.Exit(1)
	}
	go loop(l)
	<-(chan int)(nil)
}

func loop(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept (%s)", err)
			os.Exit(1)
		}
		go func() {
			defer func() {
				conn.Close()
				log.Printf("closed %s", conn.RemoteAddr())
			}()
			log.Printf("accepted %s", conn.RemoteAddr())
			for i := 0; i < 3; i++ {
				p := make([]byte, 10)
				n, _ := conn.Read(p)
				log.Printf("read from %s: buf=%s err=%v", conn.RemoteAddr(), string(p[:n]), err)
				m, err := conn.Write(p[:n])
				log.Printf("wrote to %s: buf=%s err=%v", conn.RemoteAddr(), string(p[:m]), err)
				if err != nil {
					return
				}
				conn.Write([]byte("--\n"))
			}
		}()
	}
}
