package main

func RenderRunPage() string {
	figs := A{
		"FigTwoHosts": RenderFigurePngSvg("A circuit system of two hosts (i.e. two circuit servers).", "servers", "400px"),
	}
	return RenderHtml("Running Circuit servers", Render(runBody, figs))
}

const runBody = `

<h1>Running Circuit servers</h1>

<p>Circuit servers can be started asynchronously (and in any order) 
using the command

<pre>
	circuit start -if eth0 -discover 228.8.8.8:7711
</pre>

<p>The same command is used for all instances. The <code>-if</code> option specifies the
desired network interface to bind to, while the <code>-discover</code> command 
specifies a desired IP address of a UDP multicast channel to be used for automatic
server-server discover.

The <code>-discover</code> option can be omitted by setting the environment variable
<code>CIRCUIT_DISCOVER</code> to equal the desired multicast address.

<h2>Alternative advanced server startup</h2>

<p>To run the circuit server on the first machine, pick a public IP address and port for it to
listen on, and start it like so

<pre>
	circuit start -a 10.0.0.1:11022
</pre>

<p>The circuit server will print its own circuit URL on its standard output.
It should look like this:

<pre>
	circuit://10.0.0.1:11022/78517/Q56e7a2a0d47a7b5d
</pre>

<p>Copy it. We will need it to tell the next circuit server to “join” this one
in a network, i.e. circuit.

<p>Log onto another machine and similarly start a circuit server there, as well.
This time, use the <code>-j</code> option to tell the new server to join the first one:

<pre>
	circuit start -a 10.0.0.2:11088 -j circuit://10.0.0.1:11022/78517/Q56e7a2a0d47a7b5d
</pre>

<p>You now have two mutually-aware circuit servers, running on two different
hosts in your cluster. 

{{.FigTwoHosts}}

<p>You can join any number of additional hosts to the circuit environment in a
similar fashion, even billions.  The circuit uses a modern 
<a href="http://en.wikipedia.org/wiki/Expander_graph">expander graph</a>-based algorithm for
presence awareness and ordered communication, which is genuinely distributed;
It uses communication and connectivity sparingly, hardly leaving a footprint
when idle.

        `
