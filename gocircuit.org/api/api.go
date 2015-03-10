package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderMainPage() string {
	figs := A{
		"FigHierarchy": RenderFigurePngSvg(
			"Virtual anchor hierarchy and its mapping to Go <code>Anchor</code> objects.", "hierarchy", "600px"),
	}
	return RenderHtml("Go client API", Render(mainBody, figs))
}

const mainBody = `

<h1>Go client API</h1>

<p>To use the Go client API to the circuit, start by importing the client package:
<pre>
	import "github.com/gocircuit/circuit/client"
</pre>

<h2>System abstraction</h2>

<p>The circuit organizes all of the cluster resources in an abstract hierarchichal namespace—a rooted
	tree with named edges. 
	Every node in the tree is called an <em>anchor</em> and 
		every anchor is associated with the root-to-anchor path that leads to it.
	A path identifies its anchor uniquely.
	In file system notation, paths are strings like <code>"/Xf1c8d96119cc6919/foo/bar"</code>.

<p>In addition to being a tree node in the namespace, each anchor can have none or one <em>element</em>
attached to it. An element is a logical object that manages an underlying computational resource.
There are different kinds of elements, according to their underlying resource: process, container, name server, channel, etc.

{{.FigHierarchy}}

<p>The Go client interface is organized around the anchor hierarchy abstraction. 

<p>An interface called <code>Anchor</code> represents an anchor. It provides methods
for traversing and inspecting its descendant anchors, as well as methods for creating or retrieving
the element associated it. 

<p>Circuit applications begin with a call to <code>Dial</code> which establishes connection
to a circuit cluster and returns an <code>Anchor</code> object representing the root of the
hierarchy. 

<h2>Connecting the client</h2>

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

<h2>Physical placement of anchors and panics</h2>

<p>
        `
