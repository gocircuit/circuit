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
