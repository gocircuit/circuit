The Circuit
===========

The circuit is a tiny server process which runs instances on a cluster of machines
to form a network, which enables distributed process orchestration and synchronization
from any one machine.

For a conceptual introduction to The Circuit, check out the
[GopherCon 2014 Video](http://confreaks.com/videos/3421-gophercon2014-the-go-circuit-towards-elastic-computation-with-no-failures).
Since this video was recorded, the API-via-file-system approach was abandoned
in favor of a simpler command-line tool and a Go client library.

The circuit is a tool for executing and synchronizing UNIX processes across entire clusters
by means of a command-line tool and a client library.

The circuit comes as one binary, which serves the purpose of a server
and a command-line client.

Build
-----

The Circuit comprises one small binary. It can be built for Linux and Darwin.

Given that the [Go Language](http://golang.org) compiler is installed,
you can build and install the circuit binary with one line:

	go get github.com/gocircuit/circuit/cmd/circuit

Run
---

To run the circuit agent, pick a public IP address and port for it to
listen on, and start it like so

	circuit start -a 10.0.0.7:11022

The circuit server will print its own circuit URL on its standard output.
It should look like this:

	circuit://10.0.0.7:11022/78517/Q56e7a2a0d47a7b5d

Copy it. We will need it to tell the next circuit server to “join” this one
in a network, i.e. circuit.

Log onto another machine and similarly start a circuit server there, as well.
This time, use the `-j` option to tell the new server to join the first one:

	circuit start -a 10.0.0.5:11088 -j circuit://10.0.0.7:11022/78517/Q56e7a2a0d47a7b5d

You now have two mutually-aware circuit servers, running on two different hosts in your cluster.
You can join any number of additional hosts to the circuit environment in a similar fashion,
even billions:

The circuit uses a modern [expander graph](http://en.wikipedia.org/wiki/Expander_graph)-based
algorithm for presence awareness and ordered communication, which is genuinely distributed;
It uses communication and connectivity sparingly, hardly leaving a footprint when idle.

Programming metaphor
-------

Each circuit server hosts a hierarchical namespace of control primitives,
currently _process_ or _channel_. Processes aid the execution and synchronization of
OS processes. Channels …

Elements are addressed by paths of the form

	/X317c2314a386a9db/hi/charlie

which are called _anchors_. The first part of a path is necessarily the ID of the circuit server
hosting it, while the remainder is user-specific.

Learn more
----------

The Go client for writing circuit apps is package

	github.com/gocircuit/circuit/client

The public interface of this package is self-contained. Other
packages in the circuit repo are internal.

Tutorials can be found within the client package directory

	github.com/gocircuit/circuit/client/tutorial

Additionally, the circuit binary directory contains the implementation
of the circuit tool, which is itself built using the client and is another
comprehensive example of a circuit app. It can be found in

	github.com/gocircuit/circuit/cmd/circuit

To stay up to date with new developments, documentation and articles, follow
The Circuit Project on Twitter [@gocircuit](https://twitter.com/gocircuit) or
me [@maymounkov](https://twitter.com/maymounkov).
