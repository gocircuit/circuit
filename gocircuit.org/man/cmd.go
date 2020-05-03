package man

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderCommandPage() string {
	figs := A{
		"FigClient":       RenderFigurePngSvg("Circuit client connected to a server.", "client", "500px"),
		"FigServerAnchor": RenderFigurePngSvg("Circuit servers correspond to root-level anchors.", "serveranchor", "500px"),
	}
	return RenderHtml("Command-line client", Render(commandBody, figs))
}

const commandBody = `

<h2>Using the command-line client</h2>

<p>Once the circuit servers are started, you can create, observe and control
circuit elements (i) interactively—using the circuit binary which doubles as a command-line client—as
well as (ii) programmatically—using the circuit Go client package <code>github.com/hoijui/circuit/pkg/client</code>.
In fact, the circuit command-line tool is simply a front-end for the Go client library.

<p>Clients (the tool or your own) <em>dial into</em> a circuit server in order to
interact with the entire system. All servers are equal citizens in every respect and,
in particular, any one can be used as a choice for dial-in.

{{.FigClient}}

<p>The tool (described in more detail later) is essentially a set of commands that
allow you to traverse the global hierarchical namespace of circuit elements,
and interact with them, somewhat similarly to how one uses the Zookeeper
namespace.

<p>For example, to list the entire circuit cluster anchor hierarchy, type in

<pre>
	circuit ls /
</pre>

<p>So, you might get something like this in response

<pre>
	/X88550014d4c82e4d
	/X938fe923bcdef2390
</pre>

<p>The two root-level anchors correspond to the two circuit servers.

{{.FigServerAnchor}}

<h3>Pointing the tool to your circuit cluster</h3>

<p>Before you can use the <code>circuit</code> tool, you need to tell it how to locate
one circuit server for us a <em>dial-in</em> point.

<p>There are two ways to provide the dial-in server address to the tool:

<p>1. If the circuit servers were started with the <code>-discover</code> option or the
<code>CIRCUIT_DISCOVER</code> environment variable, the command-line tool
can use the same methods for finding a circuit server. E.g.

<pre>
	circuit ls -discover 228.8.8.8:7711 /...
</pre>

<p>Or,

<pre>
	export CIRCUIT_DISCOVER=228.8.8.8:7711
	circuit ls /...
</pre>

<p>2. With the command-line option <code>-d</code>, like e.g.

<pre>
		circuit ls -d circuit://10.0.0.1:11022/78517/Q56e7a2a0d47a7b5d /
</pre>

<p>Or, equivalently, by setting the environment variable <code>CIRCUIT</code> to point to a file
whose contents is the desired dial-in address. For example, (in bash):

<pre>
		echo circuit://10.0.0.1:11022/78517/Q56e7a2a0d47a7b5d > ~/.circuit
		export CIRCUIT="~/.circuit"
		circuit ls /
</pre>

<p>A list of available tool commands is shown on the help screen

<pre>
	circuit help
</pre>

<p>A more detailed explanation of their meaning and function can be found
in the documentation of the client package, <code>github.com/gocircuit/client</code>.

`
