The Circuit
===========

The circuit is a tiny server process which runs instances on a cluster of
machines to form a network, which enables distributed process orchestration
and synchronization from any one machine.

Some of the target applications of the circuit are:

* Automatic dynamic orchestration of complex compute pipelines, as in numerical computation, for instance
* Packaging and distribution of universal distributed binaries
* Incremental automation of small and large OPS engineering workflows

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

Run the servers
---------------

To run the circuit server on the first machine, pick a public IP address and port for it to
listen on, and start it like so

	circuit start -a 10.0.0.1:11022

The circuit server will print its own circuit URL on its standard output.
It should look like this:

	circuit://10.0.0.1:11022/78517/Q56e7a2a0d47a7b5d

Copy it. We will need it to tell the next circuit server to “join” this one
in a network, i.e. circuit.

Log onto another machine and similarly start a circuit server there, as well.
This time, use the `-j` option to tell the new server to join the first one:

	circuit start -a 10.0.0.2:11088 -j circuit://10.0.0.1:11022/78517/Q56e7a2a0d47a7b5d

You now have two mutually-aware circuit servers, running on two different
hosts in your cluster. 

![A circuit system of two hosts.](https://raw.githubusercontent.com/gocircuit/circuit/master/misc/img/servers.png)

You can join any number of additional hosts to the circuit environment in a
similar fashion, even billions.  The circuit uses a modern [expander
graph](http://en.wikipedia.org/wiki/Expander_graph)-based algorithm for
presence awareness and ordered communication, which is genuinely distributed;
It uses communication and connectivity sparingly, hardly leaving a footprint
when idle.

Programming metaphor
-------

The purpose of each circuit server is to host a collection of control
primitives, called _elements_, on behalf of the user. On each server the
hosted elements are organized in a hierarchy (similarly to the file system in
Apache Zookeeper), whose nodes are called _anchors_. Anchors (akin to file
system directories) have names and each anchor can host one circuit element or
be empty.

The hierarchies of all servers are logically unified by a global circuit root
anchor, whose children are the individual circuit server hierarchies. A
typical anchor path looks like this

	/X317c2314a386a9db/hi/charlie

The first component of the path is the ID of the circuit server hosting the leaf anchor.

Except for the circuit root anchor (which does not correspond to any
particular circuit server), all other anchors can store a _process_ or a
_channel_ element, at most one, and additionally can have any number of sub-
anchors. In a way, anchors are like directories that can have any number of
subdirectories, but at most one file.

Creating and interacting with circuit elements is the mechanism through which
the  user controls and reflects on their distributed application.
This can be accomplished by means of the included Go client library, or using
the command-line tool embodied in the circuit executable itself.

Process elements are used to execute, monitor and synchronize OS-level
processes at the hosting circuit server. They allow visibility and control
over OS processes from any machine in the circuit cluster, regardless
of the physical location of the underlying OS process.

Channel elements are a synchronization primitive, similar to the channels in Go,
whose send and receive sides are accessible from any location in the
circuit cluster, while their data structure lives on the circuit server hosting
their anchor.

Use
---

Once the circuit servers are started, one can create, observe and control
circuit elements (i) interactively using the circuit binary which doubles as a command-line client,
as well as (ii) programmatically using the circuit Go client package `github.com/gocircuit/circuit/client`.
In fact, the circuit command-line tool is merely a front for a subset the Go client library,
and is a circuit client itself.

Clients (the tool or your own) connect into a circuit server and perform
operations via this server. Which server a client connects to is called
the _dial-in_ server. In general, it does not matter which server you
connect your client to. They are all equally good. And they all can control
the whole system.

![Circuit client connected to a server](https://raw.githubusercontent.com/gocircuit/circuit/master/misc/img/servers.png)

To list the entire circuit cluster anchor hierarchy, type in

	circuit ls /...

Before this command can work, however, you need to give it the address of
any one of the circuit servers as a _dial-in_ point. The choice of dial-in
server does not matter at all. All circuit servers are equally good for this job.

There are two ways to provide the dial-in address to the tool: with
the command-line option -d or by setting the environment variable `CIRCUIT`
to point to a file whose contents in the desired dial-in address.

The rest of the tool's commands can be seen by typing

	circuit help

They exactly correspond to the API of the `github.com/gocircuit/client` package,
which has a more detailed documentation and a set of tutorials.

Example: Make a process
-----------------------

Here are a few examples. To run a new process on some chosen
cluster machine, first see what machines are available:

	circuit ls /...
	---- /X88550014d4c82e4d
	---- /X938fe923bcdef2390

Run a new `ls` process:

	circuit mkproc /X88550014d4c82e4d/pippi << EOF
	{
		"Path": "/bin/ls", 
		"Args":["/"]
	}
	EOF

See what happened:

	circuit peek /X88550014d4c82e4d/pippi

Close the standard input to indicate no intention to write to it:

	cat /dev/null | circuit stdin /X88550014d4c82e4d/pippi

Read the output (note that the output won't show until you close 
the standard input first, as shown above):

	circuit stdout /X88550014d4c82e4d/pippi

Remove the process element from the anchor hierarchy

	circuit scrub /X88550014d4c82e4d/pippi

Example: Create a channel
-------------------------

Again, take a look at what servers are available:

	circuit ls /...
	---- /X88550014d4c82e4d
	---- /X938fe923bcdef2390

Pick one. Say `X88550014d4c82e4d`. Now, 
let's create a channel on `X88550014d4c82e4d`:

	circuit mkchan /X88550014d4c82e4d/this/is/charlie 3

The last argument of this line is the channel buffer capacity,
analogously to the way channels are created in Go.

Verify the channel was created:

	circuit peek /X88550014d4c82e4d/this/is/charlie

This should print out something like this:

	{
			"Cap": 3,
			"Closed": false,
			"Aborted": false,
			"NumSend": 0,
			"NumRecv": 0
	}

Sending a message to the channel is accomplished with the command

	circuit send /X88550014d4c82e4d/this/is/charlie < some_file

The contents of the message is read out from the standard input of the
command above. This command will block until a receiver is available,
unless there is free space in the channel buffer for a message.

When the command unblocks, it will send any data to the receiver.
If there is no receiver, but there is a space in the message buffer,
the command will also unblock and consume its standard input (saving
it for an eventual receiver) but only up to 32K bytes.

Receiving is accomplished with the command

	circuit recv /X88550014d4c82e4d/this/is/charlie

The received message will be produced on the standard output of 
the command above.

Be creative
-------------

The circuit allows for unusual flexibilities in process orchestration.
Take a look, for instance, at the “virus” tutorial which demonstrates
how to implement a semi-resilient self-sustained mechanism within 
a cluster. Find it in

	github.com/gocircuit/circuit/client/tutorial/3-virus

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
