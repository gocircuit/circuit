package api

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderMainPage() string {
	figs := A{
		"FigHierarchy": RenderFigurePngSvg(
			"Virtual anchor hierarchy, depicting elements attached to some of the anchors.", "hierarchy", "600px"),
		"FigResidence": RenderFigurePngSvg(
			"Except for the root, every anchor physically resides on some circuit host. "+
				"The root anchor is a logical object representing your client's connection to the cluster.",
			"residence", "630px"),
	}
	return RenderHtml("Go client API", Render(mainBody, figs))
}

const mainBody = `

<h1>Go client API</h1>

<p>To use the Go client API to the circuit, start by importing the client package:
<pre>
	import "github.com/hoijui/circuit/pkg/client"
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

<p>All circuit applications begin with a call to <code>Dial</code> (or <code>DialDiscover</code>)
which establishes connection to a circuit cluster and returns an <code>Anchor</code>
object representing the root of the hierarchy. (Programming details for connecting into a cluster
are given below.)

<h2>Anchor and element residence</h2>

<p>Every anchor—excluding the root anchor—as well as its attached element (if any) physically 
reside on some specific host in the circuit cluster. (The <code>Anchor</code> and element objects in a
Go client application are merely references to the underlying anchor and element structures.)

<p>The following illustration demonstrates how the hierarchy structure implies the physical
location of anchors and elements.

{{.FigResidence}}

<p>The root anchor (which you obtain from <code>Dial</code> or <code>DialDiscover</code>)
is special. It symbolically represents your client's connection to the circuit cluster. As such,
the root anchor resides only in your client's runtime—i.e. it is not persistent.
No elements can be attached to the root anchor.

<p>The children of the root anchor are always, by definition, server anchors.
Server anchors correspond to currently live hosts (aka servers) in your circuit cluster.
Server anchors are created and removed by the circuit system, as hosts join or leave
the circuit cluster. 

<p>Server anchors physically reside on their respective host and they have an
attached <code>Server</code> element that allows you to query various
runtime parameters of the host. <code>Server</code> elements are permanently
attached to their anchors.

<p>All anchors descendant to server anchors, and their attached elements, are created 
by the user. All such user anchors as well as the elements that might be attached to them 
reside—by definition—on the host of the server anchor that they descend from.

<h2 id="errors">Panics and errors</h2>

<p>All programmatic manipulation of a circuit client involves calling methods
of <code>Anchor</code> or element objects. As we discussed, all anchors
and elements have an implied physical place of residence (on one of the cluster hosts).

<p>In general, any method invokation might result in one of two types of errors:
<em>application errors</em> and <em>system errors</em>.

<p>Application errors are things like trying to create an element on anchor that
already has one, or trying to start a process using a missing binary, for instance.
Such errors will be returned in the form of Go <code>error</code> return values
of the respective method.

<p>Independently of application errors, every invokation of an anchor or
element method may fail if the underlying object is physically unreachable.
Anchors residing on a dead host are unreachable and so are their elements,
for example. Such errors are treated in a separate category of system errors
and they are reported as panics. In particular, if a host is unreachable, 
all anchors descendant to and including its server anchor will cause panics when used.

<p>By design, any anchor or element method invokation will result in
a panic, if a system error occurs. We uniformly report system errors as 
panics in order to separate them semantically from application errors.
But also because they have asynchronous nature and because they usually
result in a very different way of being handled by the application programmer.

<p>That said, such panic conditions are not critical. These panics merely
indicate that the host where an anchor or element physically resides is currently
unreachable. The underlying host can be unreachable either if dead or as the
result of a complete network partition (partial partitions do not affect the system).

<p>An anchor or element object that produces a panic remains in a valid
state after the panic and it can be re-used. If the underlying resource is still
unreachable, another panic will be produced. But if the system has
recovered from a network partition and the underlying resource is reachable
again, follow on method calls will succeed.

<h3>Connection panics</h3>

<p>Panics in any method invocation can also be caused if the client's connection
to a circuit server is lost. This type of panic is permanent, as the circuit client
does not attempt automatic reconnection to the circuit cluster.

<p>There is a way to distinguish between host-only panics and permanent client
connection panics. After catching a panic anywhere, the user application can
simply call the root anchor's <code>View</code> method (which lists the contents of
the anchor). If this call also results in a panic, this is an indication that the client
connection has been lost altogether. 

        `
