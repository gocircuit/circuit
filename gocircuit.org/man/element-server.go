package man

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderElementServerPage() string {
	return RenderHtml("Circuit channel element", Render(elementServerBody, nil))
}

const elementServerBody = `

<h2>Example: Using server elements</h2>

<p>Server elements are interfaces to the circuit daemons that comprise the cluster.
They provide information about the daemon, as well as the ability to merge two
disconnected circuit clusters.

<p>Let us go through an example. First, take a look at what servers are available:

<pre>
	circuit ls /...
	/X88550014d4c82e4d
	/X938fe923bcdef2390
</pre>

<p>Pick <code>X88550014d4c82e4d</code>, for instance. 
The peek command will give you basic information about the daemon at this
path, including its circuit address and the time it was started:

<pre>
	# circuit peek /X88550014d4c82e4d
	{
	    "Addr": "circuit://127.0.0.1:23111/5347/Q7315aad801e4994d",
	    "Joined": "2015-03-18T13:35:08.217020207-04:00"
	}
</pre>

<p>The stack trace command will print the Go stack trace of the running circuit daemon:
<pre>
	# circuit stk /X88550014d4c82e4d
</pre>

<p>Finally, using the join command, we can instruct the circuit daemon underlying the path 
<code>/X88550014d4c82e4d</code> to merge the cluster that it is a part of with the
cluster that contains a given circuit daemon, specified by its circuit address.
(A circuit address provides a way of connecting to a circuit daemon directly.)

<pre>
	# circuit join /X88550014d4c82e4d circuit://127.0.0.1:41222/5650/Q4e16779fe039ecf3
</pre>
<p>This command will instruct the circuit daemon underlying <code>/X88550014d4c82e4d</code>
to perform networking operations that will result in joining this circuit cluster with the 
circuit cluster that the target address <code>circuit://127.0.0.1:41222/5650/Q4e16779fe039ecf3</code> is
a part of. If the target is already a member of this cluster, no change will occur.

        `
