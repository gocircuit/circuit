nix:
	go build -o $(GOPATH)/bin/circuit cmd/circuit/main.go cmd/circuit/server.go cmd/circuit/hmac.go \
		cmd/circuit/chan.go cmd/circuit/dns.go cmd/circuit/joinleave.go \
		cmd/circuit/ls.go cmd/circuit/wait.go cmd/circuit/util.go \
		cmd/circuit/start.go cmd/circuit/procdkr.go cmd/circuit/peek.go \
		cmd/circuit/load.go cmd/circuit/std.go cmd/circuit/recv.go

clean:
	rm $(GOPATH)/bin/circuit
