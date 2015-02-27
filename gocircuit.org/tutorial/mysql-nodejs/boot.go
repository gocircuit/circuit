package mysql_nodejs

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderMysqlNodejsBoot() string {
	return RenderHtml("Starting a MySQL and node.js stack using a circuit app", Render(bootBody, nil))
}

const bootBody = `

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
