
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
For instance, `/circuit` is a good choice.

To run the circuit tool on local port `11022`, use

	circuit -a :11022 -m /circuit

Among a few other things, the circuit tool will printout its own circuit URL.
It should look like this:

	…
	circuit://[::]:11022/78139/R21b66be46e9ba3e8
	…

Copy it. We will need it to instruct the next circuit tool to join this one.
