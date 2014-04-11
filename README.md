

The Circuit is a tool for executing and synchronizing OS processes across entire clusters
by means of a file-system interface.


Build
-----

The Circuit comprises one small binary. It can be built for Linux and Darwin.

Given that the [Go Language](http://golang.org) compiler is installed,
you can build and install the circuit binary with one line:

	go get github.com/gocircuit/circuit/cmd/circuit

Run
---

Prepare a local directory that can be FUSE-mounted by your user. 
For instance, `/circuit` is a good choice.

To run the circuit agent, pick a public IP address and port for it to
listen on, and start it like so

	circuit -a 10.20.30.7:11022 -m /circuit

Among a few other things, the circuit agent will print its own circuit URL.
It should look like this:

	…
	circuit://10.20.30.7:11022/78517/R56e7a2a0d47a7b5d
	…

Copy it. We will need it to tell the next circuit agent to join this one.

Log onto another machine and similarly start a circuit agent there, as well.
This time, use the `-j` option to tell the new agent to join the
circuit of the first one:

	circuit -a 10.20.30.5:11088 -m /circuit -j circuit://10.20.30.7:11022/78517/R56e7a2a0d47a7b5d

You now have two mutually-aware circuit agents, running on two different hosts in your cluster.
You can join any number of additional hosts to the circuit environment in a similar fashion.

Explore
-------

On any host with a running circuit agent, go to the local circuit mount directory

	cd /circuit
	ls

Each of its subdirectories corresponds to a live circuit agent. Navigate into
any one of them and explore the file system. Each directory is equipped with a
`help` file to guide you.
