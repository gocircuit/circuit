# Circuit

[![Build Status](https://drone.io/github.com/gocircuit/circuit/status.png)](https://drone.io/github.com/gocircuit/circuit/latest) [![GoDoc](https://godoc.org/github.com/gocircuit/circuit/client?status.png)](https://godoc.org/github.com/gocircuit/circuit/client)

![Engineering role separation.](https://raw.githubusercontent.com/gocircuit/circuit/master/misc/img/3.png)

**The CIRCUIT is a new way of thinking. It is deceptively similar to existing software,
while being quite different.**

Circuit is a programmable platform-as-a-service (PaaS) and/or Infrastructure-as-a-Service (IaaS), 
for management, discovery, synchronization and orchestration of services and 
hosts comprising cloud applications. 

Circuit was designed to enable clear, accountable and safe interface between the human engineering
roles in a technology enterprise, ultimately increasing productivity. Engineering role separation
in a typical circuit-based architecture is illustrated above.

![A circuit-managed cloud.](https://raw.githubusercontent.com/gocircuit/circuit/master/misc/img/header.png)

Users of circuit are

* Operations engineers, who sustain cloud applications at host, process and network level
* Data scientists, who develop distributed compute pipelines by linking together and distributing third-party utilities
* Manufacturers of distributed software, who wish to codify installation and maintenance procedures in a standardized
fashion instead of communicating them through documentation (viz. MySQL)

A few technical features of circuit:

* Single- and multi-datacenter out-of-the-box
* Authentication and security of system traffic
* TCP- and UDP multicast-based (similarly to [mDNS](http://en.wikipedia.org/wiki/Multicast_DNS)) internal node discovery
* Zero-configuration/blind-restart for ease on the host provisioning side
* Global, point-of-view consistent key/value space, where keys are hierarchical paths and 
values are control objects for data, processes, synchronizations, and so on.
* Consistency guarantees at ultra-high churn rates of physical hosts
* Command-line and programmatic access to the API
* Integration with Docker

In a typical circuit scenario: 

* Provisioning engineers ensure newly provisioned
machines start the zero-configuration circuit server as a daemon on startup.
* Operations engineers program start-up as well as dynamic-response 
behavior via command-line tool or language bindings.

Adoption considerations:

* Small footprint: Circuit daemons leave almost no communication and memory footprint when left idle.
This makes circuit ideal for incremental adoption alongside pre-existing architectures
* Immediate impact: Even short and simple circuit scripts save manual time going forward
* Knowledge accounting: Circuit scripts can replace textual post-mortem reports with
executable discovery, diagnosis, action and inaction recipes.
* Circuit servers log all operations in their execution orders, enabling maximum visibility
during post-factum debugging and analysis.

Programming environment:

* Circuit programs (sequences of invocations of the circuit API) are not
declarative (as in Puppet, Chef, etc.). They are effectively imperative
programs in the [CSP](http://en.wikipedia.org/wiki/Communicating_sequential_processes) concurrency model,
which allows engineers to encode complex dynamic response behavior, spanning multiple data centers.

Find comparisons to other technologies—like Zookeeper, etcd, CoreOS, raft, Consul, Puppet, Chef, and 
so forth—in [the wiki](https://github.com/gocircuit/circuit/wiki).

## Incomparable but related works

* [Google Cloud Platform](https://github.com/GoogleCloudPlatform) and [Kubernetes](https://github.com/GoogleCloudPlatform/kubernetes)
* [Hashicorp Consul](https://github.com/hashicorp/consul)
* [Puppet](http://puppetlabs.com/)
* [Chef](http://www.getchef.com/solutions/devops/)
* [Ansible](http://www.ansible.com/home)
* [CoreOS](https://coreos.com/)
* [Apache Zookeeper](http://zookeeper.apache.org/)
* [Polymer (Google)](http://www.polymer-project.org/)
* [Julia (MIT)](http://julialang.org/)
* and so on.

None of these related products sees the cluster as a closed system. In this way,
the circuit is different than all. This is explained in precise terms in the next section.

## Integration

The circuit is a tiny server process which runs instances on a cluster of
machines to form an efficient, churn-resilient network, which enables distributed process orchestration
and synchronization from any one machine.

Some of the target applications of the circuit are:

* Automatic dynamic orchestration of complex compute pipelines, as in numerical computation, for instance
* Packaging and distribution of universal distributed binaries that can self-organize into complex cloud apps
* Incremental automation of small and large OPS engineering workflows

Dive straight into it with the [Quick Start](https://docs.google.com/presentation/d/1x6ZqTg6AeVHMGh1ug1oqGZ9kYFLoorjRhDlTs7ytj3o/edit?usp=sharing) slide deck.

For a conceptual introduction to The Circuit, check out the
[GopherCon 2014 Video](http://confreaks.com/videos/3421-gophercon2014-the-go-circuit-towards-elastic-computation-with-no-failures).
Since this video was recorded, the API-via-file-system approach was abandoned
in favor of a simpler command-line tool and a Go client library.

Also take a look at the [faux animated illustration](https://docs.google.com/presentation/d/1nazPJAmYeIvJam9oiA-6euEz3Bo9flRLqsgRU7QAlr0/edit?usp=sharing)
of the [Advanced Tutorial: Watchbot with a back channel](https://github.com/gocircuit/circuit/tree/master/tutorial/watchbot-with-chan).

The circuit is a tool for executing and synchronizing UNIX processes across entire clusters
by means of a command-line tool and a client library.

The circuit comes as one binary, which serves the purpose of a server
and a command-line client.

## Build ##

The Circuit comprises one small binary. It can be built for Linux and Darwin.

Given that the [Go Language](http://golang.org) compiler is [installed](http://golang.org/doc/install),
you can build and install the circuit binary with one line:

	go get github.com/gocircuit/circuit/cmd/circuit

## Run the servers

Circuit servers can be started asynchronously (and in any order) 
using the command

	circuit start -if eth0 -discover 228.8.8.8:7711

The same command is used for all instances. The `-if` option specifies the
desired network interface to bind to, while the `-discover` command 
specifies a desired IP address of a UDP multicast channel to be used for automatic
server-server discover.

The `-discover` option can be omitted by setting the environment variable
`CIRCUIT_DISCOVER` to equal the desired multicast address.

### Alternative advanced server startup 

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

## Programming metaphor ##

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

## Use ##

Once the circuit servers are started, you can create, observe and control
circuit elements (i) interactively—using the circuit binary which doubles as a command-line client—as
well as (ii) programmatically—using the circuit Go client package `github.com/gocircuit/circuit/client`.
In fact, the circuit command-line tool is simply a front-end for the Go client library.

Clients (the tool or your own) _dial into_ a circuit server in order to
interact with the entire system. All servers are equal citizens in every respect and,
in particular, any one can be used as a choice for dial-in.

![Circuit client connected to a server](https://raw.githubusercontent.com/gocircuit/circuit/master/misc/img/client.png)

The tool (described in more detail later) is essentially a set of commands that
allow you to traverse the global hierarchical namespace of circuit elements,
and interact with them, somewhat similarly to how one uses the Zookeeper
namespace.

For example, to list the entire circuit cluster anchor hierarchy, type in

	circuit ls /

So, you might get something like this in response

	/X88550014d4c82e4d
	/X938fe923bcdef2390

The two root-level anchors correspond to the two circuit servers.

![Circuit servers correspond to root-level anchors](https://raw.githubusercontent.com/gocircuit/circuit/master/misc/img/serveranchor.png)

### Pointing the tool to your circuit cluster ###

Before you can use the `circuit` tool, you need to tell it how to locate
one circuit server for us a _dial-in_ point.

There are two ways to provide the dial-in server address to the tool:

1. If the circuit servers were started with the `-discover` option or the
`CIRCUIT_DISCOVER` environment variable, the command-line tool
can use the same methods for finding a circuit server. E.g.

	circuit ls -discover 228.8.8.8:7711 /...

Or,

	export CIRCUIT_DISCOVER=228.8.8.8:7711
	circuit ls /...

2. With the command-line option `-d`, like e.g.

		circuit ls -d circuit://10.0.0.1:11022/78517/Q56e7a2a0d47a7b5d /

Or, equivalently, by setting the environment variable `CIRCUIT` to point to a file
whose contents is the desired dial-in address. For example, (in bash):

		echo circuit://10.0.0.1:11022/78517/Q56e7a2a0d47a7b5d > ~/.circuit
		export CIRCUIT="~/.circuit"
		circuit ls /

A list of available tool commands is shown on the help screen

	circuit help

A more detailed explanation of their meaning and function can be found
in the documentation of the client package, `github.com/gocircuit/client`.

### Example: Make a process ###

Here are a few examples. To run a new process on some chosen
cluster machine, first see what machines are available:

	circuit ls /...
	/X88550014d4c82e4d
	/X938fe923bcdef2390

Run a new `ls` process:

	circuit mkproc /X88550014d4c82e4d/pippi << EOF
	{
		"Path": "/bin/ls", 
		"Args":["/"]
	}
	EOF

![Process elements execute OS processes on behalf of the user](https://raw.githubusercontent.com/gocircuit/circuit/master/misc/img/mkproc.png)

See what happened:

	circuit peek /X88550014d4c82e4d/pippi

Close the standard input to indicate no intention to write to it:

	cat /dev/null | circuit stdin /X88550014d4c82e4d/pippi

Read the output (note that the output won't show until you close 
the standard input first, as shown above):

	circuit stdout /X88550014d4c82e4d/pippi

Remove the process element from the anchor hierarchy

	circuit scrub /X88550014d4c82e4d/pippi

### Example: Make a docker container ###

Much like for the case of OS processes, the circuit can create, manage and synchronize [Docker](http://www.docker.com) containers,
and attach the corresponding _docker elements_ to a path in the anchor file system.

To allow creation of docker elements, any individual server must be started with the `-docker` switch. For instance:

	circuit start -if eth0 -discover 228.8.8.8:7711 -docker

To create and execute a new docker container, using the tool:

	circuit mkdkr /X88550014d4c82e4d/docky << EOF
	{
		"Image": "ubuntu",
		"Memory": 1000000000,
		"CpuShares": 3,
		"Lxc": ["lxc.cgroup.cpuset.cpus = 0,1"],
		"Volume": ["/webapp", "/src/webapp:/opt/webapp:ro"],
		"Dir": "/",
		"Entry": "",
		"Env": ["PATH=/usr/bin"],
		"Path": "/bin/ls",
		"Args": ["/"],
	}
	EOF

Most of these fields can be omitted analogously to their command-line option counterparts 
of the `docker` command-line tool.

![Docker elements are like processes](https://raw.githubusercontent.com/gocircuit/circuit/master/misc/img/mkdkr.png)

The remaining docker element commands are identical to those for processes:
`stdin`, `stdout`, `stderr`, `peek` and `wait`. In one exception, `peek` will return
a detailed description of the container, derived from `docker inspect`. 

### Example: Create a channel ###

Again, take a look at what servers are available:

	circuit ls /...
	/X88550014d4c82e4d
	/X938fe923bcdef2390

Pick one. Say `X88550014d4c82e4d`. Now, 
let's create a channel on `X88550014d4c82e4d`:

	circuit mkchan /X88550014d4c82e4d/this/is/charlie 3

The last argument of this line is the channel buffer capacity,
analogously to the way channels are created in Go.

![Channel elements reside in the memory of a circuit server](https://raw.githubusercontent.com/gocircuit/circuit/master/misc/img/mkchan.png)

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

### Example: Make a DNS server element ###

Circuit allows you to create and dynamically configure one or more DNS
server elements on any circuit server.

As before, pick an available circuit server, say `X88550014d4c82e4d`.
Create a new DNS server element, like so

	circuit mkdns /X88550014d4c82e4d/mydns

This will start a new DNS server on the host of the given circuit server,
binding it to an available port. Alternatively, you can supply an
IP address argument specifying the bind address, as in

	circuit mkdns /X88550014d4c82e4d/mydns 127.0.0.1:7711

Either way, you can always retrieve the address on which the DNS server
is listening by peeking into the corresponding circuit element:

	circuit peek /X88550014d4c82e4d/mydns

This command will produce an output similar to this

	{
	    "Address": "127.0.0.1:7711",
	    "Records": {}
	}

Once the DNS server element has been created, you can add resource records
to it, one at a time, using

	circuit set /X88550014d4c82e4d/mydns "miek.nl. 3600 IN MX 10 mx.miek.nl."

Resource records use the canonical syntax, described in various RFCs.
You can find a list of such RFCs as well as examples in the DNS Go library
that underlies our implementation: `github.com/miekg/dns/dns.go`

All records, associated with a given name can be removed with a single command:

	circuit unset /X88550014d4c82e4d/mydns miek.nl.

The current set of active records can be retrieved by peeking into the element:

	circuit peek /X88550014d4c82e4d/mydns

Assuming that a name has multiple records associated with it, peeking would produce
an output similar to this one:

	{
		"Address": "127.0.0.1:7711",
		"Records": {
			"miek.nl.": [
				"miek.nl.\t3600\tIN\tMX\t10 mx.miek.nl.",
				"miek.nl.\t3600\tIN\tMX\t20 mx2.miek.nl."
			]
		}
	}

### Example: Listen on server join and leave announcements ###

The circuit provides two special element types `@join` and `@leave`, 
called _subscriptions_. Their job is to notify you when new
circuit servers join the systems or others leave it.
Both of them behave like receive-only channels.

	circuit mk@join /X88550014d4c82e4d/watch/join
	circuit mk@leave /X88550014d4c82e4d/watch/leave

The join subscription delivers a new message each time a cicruit server joins
the system. The received message holds the anchor path of the new server.

	circuit recv /X88550014d4c82e4d/watch/join

Similarly, the leave subscription delivers a new message each time
a circuit server disappears from the system.

	circuit peek /X88550014d4c82e4d/watch/join

## Be creative ##

The circuit allows for unusual flexibilities in process orchestration.
Take a look, for instance, at the two “watchbot” tutorials which demonstrate
how to implement a semi-resilient self-sustained mechanism within 
a cluster. Find the simpler one here

	https://github.com/gocircuit/circuit/tree/master/tutorial/watchbot

And the more elaborate one, which demonstrate use of channels, here

	https://github.com/gocircuit/circuit/tree/master/tutorial/watchbot-with-chan

## Security ##

By default, circuit servers and clients communicate over plaintext TCP.
A HMAC-based symmetric authentication, followed by an asymmetric
RC4 stream cipher is supported.

To enable encryption, use the `-hmac` command-line option to point
the circuit executable to a file containing the private key for your circuit.
For instance:

	circuit start -a 10.0.0.1 -hmac .hmac

Or, if you are invoking the tool:

	circuit ls -hmac .hmac /...

Alternatively, you can set the environment `CIRCUIT_HMAC` to
point to the private key file.

To generate a new private key for your circuit, use the command

	circuit keygen

## Networking ##

From a networking and protocol standpoint, circuit servers and
clients are peers: All communications (server-server and server-client)
use a common RPC framework which often entails a server
being able to reverse-dial into a client.

For this reason, circuit clients (the circuit tool or your apps) CANNOT
be behind a firewall with respect to the servers they are dialing into.

# Learn more #

Ask questions to [The Circuit User Group](https://groups.google.com/forum/#!forum/gocircuit-user).

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

## Donate with a tweet ##

Please, tweet about the circuit (mention [@gocircuit](https://twitter.com/gocircuit)).
Tweets are a much appreciated donation and they help us and our sponsors
gauge the interest in this unconventional idea.

## Sponsors and credits

* [DARPA XDATA](http://www.darpa.mil/Our_Work/I2O/Programs/XDATA.aspx) initiative, 2012–2014
* [Data Tactics Corp](http://www.data-tactics.com/), 2012-2014
* [L3](http://www.l-3com.com/), 2014
* [Tumblr, Inc.](http://tumblr.com), 2012

## Featured

* [DARPA](http://www.darpa.mil) [Open Catalog](http://www.darpa.mil/opencatalog/), Arlington, VA, 2014
* [GOPHERCON 2014](http://confreaks.com/videos/3421-gophercon2014-the-go-circuit-towards-elastic-computation-with-no-failures) Denver, CO
* [STRANGELOOP 2013](http://blog.gocircuit.org/strangeloop-2013), St. Louis, MO

## Bibliography

* [Old Circuit website, original white papers and design documents](http://gocircuit-org.appspot.com)
