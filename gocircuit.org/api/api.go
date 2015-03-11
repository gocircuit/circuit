package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
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

<h2>Panics and errors</h2>

        `
