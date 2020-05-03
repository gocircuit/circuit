package main

import (
	"fmt"
	"net"
	// "time"

	. "github.com/hoijui/circuit/github.com/miekg/dns"
)

func Filter(w ResponseWriter, req *Msg) {
	m := new(Msg)
	m.SetReply(req)
	q := m.Question[0].Name // question, e.g. "miek.nl."
	println("querying:", q)

	m.Answer = make([]RR, 1)
	m.Answer[0] = &A{
		Hdr: RR_Header{
			Name:   q,
			Rrtype: TypeA,
			Class:  ClassINET,
			Ttl:    0,
		},
		A: net.IP{127, 0, 0, 1},
	}
	w.WriteMsg(m)
}

func RunLocalUDPServer(laddr string) (*Server, string, error) {
	pc, err := net.ListenPacket("udp", laddr)
	if err != nil {
		return nil, "", err
	}
	mux := NewServeMux()
	mux.HandleFunc(".", Filter)
	server := &Server{
		PacketConn: pc,
		Handler:    mux,
	}
	go func() {
		server.ActivateAndServe()
		pc.Close()
	}()
	// go func() {
	// 	time.Sleep(5*time.Second)
	// 	println("shutting...")
	// 	server.Shutdown()
	// }()
	return server, pc.LocalAddr().String(), nil
}

func RunLocalTCPServer(laddr string) (*Server, string, error) {
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		return nil, "", err
	}
	mux := NewServeMux()
	mux.HandleFunc(".", Filter)
	server := &Server{
		Listener: l,
		Handler:  mux,
	}
	go func() {
		server.ActivateAndServe()
		l.Close()
	}()
	return server, l.Addr().String(), nil
}

func main() {
	_, addr, err := RunLocalUDPServer("127.0.0.1:61222")
	fmt.Printf("addr=%v err=%v\n", addr, err)
	select {}
}
