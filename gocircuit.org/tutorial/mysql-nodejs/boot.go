package mysql_nodejs

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderBoot() string {
	return RenderHtml("Boot the circuit cluster", Render(bootBody, nil))
}

const bootBody = `

<h1>Boot the circuit cluster</h1>

<p>Booting the cluster involves starting your hosts, running the circuit on
each of them, and making sure that all circuit daemons are connected.

<p>For the tutorial we need two hosts ideally, but you can start any
desired number 1, 2, 3, …

<h3>Start the host instance</h3>

<p>Begin by starting a new EC2 instance, using the image 
<a href="tutorial-mysql-nodejs-image.html">we created</a>.

<p>This and all subsequent instances will belong to the same virtual
private network inside EC2. EC2 will, as default, give two IP addresses
to each host — one private, one public. 

<p>Host instances on your network will be able to connect to all
ports on other host instances, using their private IP address.
Hosts outside EC2 can connect to a restricted set of ports on
the public IP addresses. 

<p>When configuring the first host instance before launch, 
in addition to the requisite SSH port, make sure to leave open
another TCP port (say 11022) if you would like to be able
to connect to the circuit directly from your notebook.
These configurations are accomplished in the “security group”
section on EC2.

<h3>Start the circuit server</h3>

<p>Once the host instance is running, connect into it using SSH.

<p>Discover the private address of the EC2 host instance, and save it into a variable:
<pre> 
	# ip_address=` + "`" + `curl http://169.254.169.254/latest/meta-data/local-ipv4` + "`" + `
</pre> 

<p>Start the circuit server, instructing it to listen on the private IP address of this host,
on port 11022. (The port choice is arbitrary.)
<pre> 
	# circuit start -a ${ip_address}:11022 1> /var/circuit/address 2> /var/circuit/log &
</pre>

<p>When the server starts it will print its circuit address (a URL-type string that
looks like <code>circuit://…</code>) to standard output and we
save this into the host-local file <code>/var/circuit/address</code> for future use.
The server logs all commands and other events that happen to it to standard error.
Respectively, we redirect it to a host-local log file named <code>/var/circuit/log</code>.

<p>This start-up procedure can be summarized in a shell script, which you can
locate in the source repo at

<pre>
	$GOPATH/src/github.com/hoijui/circuit/tutorial/ec2/start-first-server.sh
</pre>

<p>Note that we are starting the circuit server as a singleton server, without
specifying an automatic method for discovering other servers. This is because
Amazon EC2 does not support Multicast UDP, which is needed for automatic
server-server discovery.

<p>Launch the remaining hosts and their respective circuit servers in the manner
described here. Naturally if you are building a production system, you would
preconfigure the host image to start the circuit server automatically upon booting.

<p>Suppose you launch a total of three hosts using this method. You now have three hosts
with circuit server running on each. At this stage each circuit server is
a singleton member in a network of one. Our next step will be to join
all three of them in a single network of three circuit servers.

<h3>Connect circuit servers into a single cluster</h3>

<p>Obtain the circuit address of each of the running circuit servers. On any
given host you can accomplish this by saying

<pre>
	# cat /var/circuit/address
</pre>

<p>You should expect to see a URL-type string that looks like this in return:

<pre>
	circuit://174.51.10.12:11022/12477/Q15d828c92f4c90ff
</pre>

<p>This URL uniquely identifies the circuit server process that created it, as well as
describes how to connect to the server.

<p>Next we are going to use the circuit client to log into one of the servers, say <code>host1</code>,
and instruct it to join the network of each of the other two servers.
Note that joining is a symmetric operation. If server A is instructed to join server B,
the result is that the networks that A and B are part of are merged into one. (And
if A and B were already part of the same network, no change occurs.)

<p>Log into <code>host1</code> and, for convenience, save the circuit address of the local
server into a shell variable <code>H1</code>:

<pre>
	host1# H1=$(cat /var/circuit/address)
</pre>

<p>Use the circuit client to connect into the local circuit server and list all servers
that it sees in its own network:

<pre>
	host1# circuit ls -d $H1 /
</pre>

<p>You should see a single member — the circuit server running on <code>host1</code> itself — 
listed by its name in the circuit's virtual file system.
The output should look something like:

<pre>
	/Xfea8b5b798f2fc09
</pre>

<p>Place the circuit addresses of the other two hosts, whatever they might be in your case,
in the variables <code>H2</code> and <code>H3</code>. Using the circuit client again,
instruct the local circuit server (the on <code>host1</code>) to join the networks of the other two:

<pre>
	host1# circuit join -d $H1 /Xfea8b5b798f2fc09 $H1
	host1# circuit join -d $H1 /Xfea8b5b798f2fc09 $H2
</pre>

<p>Let us break down these command lines. The part <code>-d $H1</code>
tells the client how to connect into the local server (the circuit
address <code>$H1</code> contains the host and port of the running server). 
Then it instructs the server with virtual name <code>/Xfea8b5b798f2fc09</code> —
which happens to be the local server we are connected into — to join
its network into that of <code>$H1</code> and <code>$H2</code> respectively.

<p>As a result, all three servers will become part of the same network,
and you can verify this by listing the members of the network again.

<pre>
	host1# circuit ls -d $H1 /
</pre>

<p>This time you should see three entries, along the lines of:

<pre>
	/X2987b5b023f2f988
	/Xca2b345798112c09
	/Xfea8b5b798f2fc09
</pre>

<p>If you were to log into any of the other hosts and use the circuit client from there,
say from the second host:

<pre>
	host2# circuit ls -d $H2 /
</pre>

<p>You should see the exact same list of three members.

<p>At this stage the circuit cluster is connected and ready to be used.
You should not have to restart and rejoin circuit servers unless a host 
dies and/or you are adding a new one to the cluster.

        `
