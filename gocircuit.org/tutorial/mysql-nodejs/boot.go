package mysql_nodejs

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
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

<p>When the server starts it will print its circuit address to standard output and we
save this into the host-local file <code>/var/circuit/address</code> for future use.
The server logs all commands and other events that happen to it to standard error.
Respectively, we redirect it to a host-local log file named <code>/var/circuit/log</code>.

<p>This start-up procedure can be summarized in a shell script, which you can
locate in the source repo at

<pre>
	$GOPATH/src/github.com/gocircuit/circuit/tutorial/ec2/start-first-server.sh
</pre>

<p>Note that we are starting the circuit server as a singleton server, without
specifying an automatic method for discovering other servers. This is because
Amazon EC2 does not support Multicast UDP, which is needed for automatic
server-server discovery.

<p>Instead, 

<h3>xx</h3>

<p>Add the following shell script <code>/usr/local/bin/start-first-server.sh</code>:

<pre>
	#!/bin/sh
	# Save the EC2 private IP address of this host to a variable.
	ip_address=` + "`" + `ifconfig eth0 | awk '/inet addr/ {split($2, a, ":"); print a[2] }'` + "`" + `
	# Start the circuit server
	/usr/local/bin/circuit start -a ${ip_address}:11022 1> /var/circuit/address 2> /var/circuit/log &
</pre>

<p>Add the following shell script <code>/usr/local/bin/start-joining-server.sh</code>:

<pre>
	#!/bin/sh
	# Save the EC2 private IP address of this host to a variable.
	ip_address=` + "`" + `ifconfig eth0 | awk '/inet addr/ {split($2, a, ":"); print a[2] }'` + "`" + `
	# Start the circuit server
	/usr/local/bin/circuit start -a ${ip_address}:11022 -j $1 1> /var/circuit/address 2> /var/circuit/log &
</pre>


        `
