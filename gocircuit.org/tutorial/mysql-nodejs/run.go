package mysql_nodejs

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderRun() string {
	return RenderHtml("Run the app on the cluster", Render(runBody, nil))
}

const runBody = `
<h1>Run the app on the cluster</h1>

<p>Here we get to run the circuit program that we wrote in the previous section.
Log into any one of the EC2 instances that are part of your circuit cluster.

<p>First, build and install the circuit app, which can be found within the circuit repo:
<pre>
	$ go install github.com/hoijui/circuit/tutorial/nodejs-using-mysql/start-app
</pre>

<p>This will place the resulting executable in <code>$GOPATH/bin</code>.

<p>And finally we can execute the circuit app, instructing it to connect (as a client)
to the circuit server running on the host we are currently on:
<pre>
	$ $GOPATH/bin/start-app -addr $(cat /var/circuit/address)
</pre>
<p>If successful, the app will print out the addresses of the MySQL server and
the Node.js service and will exit.

<p>You should also be able to see the Node.js process element using the circuit command-line tool:
<pre>
	$ circuit ls -d $(cat /var/circuit/address) -l /...
</pre>

<p>This concludes the tutorial.

        `
