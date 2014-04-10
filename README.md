
The Circuit is a simple language-agnostic tool for programming process-level concurrency across a cluster


Build and install
-----------------

The Circuit comprises a single, small (about 6MB)  binary.
Assuming the [Go Language](http://golang.org) compiler is installed,
you can build and install the circuit binary with one line:

	go get github.com/gocircuit/circuit/cmd/circuit

Run
---

Prepare a local directory that can be FUSE-mounted by your user. 
For instance, _/circuit_ is a good choice.

To run the circuit tool on local port _11022_, use

	circuit -a :11022 -m /circuit

