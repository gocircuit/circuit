package main

func RenderMetaphorPage() string {
	return RenderHtml("Programming metaphor", Render(metaphorBody, nil))
}

const metaphorBody = `

<h1>Programming metaphor </h1>

<p>The purpose of each circuit server is to host a collection of control
primitives, called <em>elements</em>, on behalf of the user. On each server the
hosted elements are organized in a hierarchy (similarly to the file system in
Apache Zookeeper), whose nodes are called <em>anchors</em>. Anchors (akin to file
system directories) have names and each anchor can host one circuit element or
be empty.

<p>The hierarchies of all servers are logically unified by a global circuit root
anchor, whose children are the individual circuit server hierarchies. A
typical anchor path looks like this

<pre>
	/X317c2314a386a9db/hi/charlie
</pre>

<p>The first component of the path is the ID of the circuit server hosting the leaf anchor.

<p>Except for the circuit root anchor (which does not correspond to any
particular circuit server), all other anchors can store a <em>process</em> or a
<em>channel</em> element, at most one, and additionally can have any number of sub-
anchors. In a way, anchors are like directories that can have any number of
subdirectories, but at most one file.

<p>Creating and interacting with circuit elements is the mechanism through which
the  user controls and reflects on their distributed application.
This can be accomplished by means of the included Go client library, or using
the command-line tool embodied in the circuit executable itself.

<p>Process elements are used to execute, monitor and synchronize OS-level
processes at the hosting circuit server. They allow visibility and control
over OS processes from any machine in the circuit cluster, regardless
of the physical location of the underlying OS process.

<p>Channel elements are a synchronization primitive, similar to the channels in Go,
whose send and receive sides are accessible from any location in the
circuit cluster, while their data structure lives on the circuit server hosting
their anchor.

        `
