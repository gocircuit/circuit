package api

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderConnectPage() string {
	figs := A{
		"FigHierarchy": RenderFigurePngSvg(
			"Virtual anchor hierarchy, depicting elements attached to some of the anchors.", "hierarchy", "600px"),
	}
	return RenderHtml("Connecting to a circuit cluster", Render(connectBody, figs))
}

const connectBody = `

<h1>Connecting to a circuit cluster</h1>

<p>To use the Go client API to the circuit, start by importing the client package:
<pre>
	import "github.com/hoijui/circuit/client"
</pre>

<p>The first step of every circuit client application is connecting to a circuit cluster.
There are two alternative methods for doing so.

<h3>Connecting to a specific server</h3>

<p>If you want to connect 
the client to a specific circuit server with a known circuit address, then use
<pre>
Dial(addr string, authkey []byte) *Client
</pre>
<p>Argument <code>addr</code> specifies the circuit server address (a string of the form <code>circuit://…</code>),
	whereas <code>authkey</code> should equal the authentication key for this cluster, or <code>nil</code> if
	the cluster does not use encryption and authentication.

<p><code>Dial</code> blocks until the connection is established and on success returns a client
object, which implements the <code>Anchor</code> interface and represents the root of the 
anchor hierarchy.

<p>If <code>Dial</code> is to fail, it will report an error by panicing.

<h3>Connecting by discovering a server</h3>

<p>Alternatively, if the circuit cluster supports (multicast-based) discovery, 
you could use <code>DialDiscover</code> to first discover a random
circuit server from the cluster and then connect to it:
<pre>
DialDiscover(multicast string, authkey []byte) *Client
</pre>

<p>The argument <code>multicast</code> must equal the multicast discovery address for the
circuit cluster.

        `
