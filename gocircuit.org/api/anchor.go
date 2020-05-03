package api

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderAnchorPage() string {
	figs := A{
		"FigHierarchy": RenderFigurePngSvg(
			"Virtual anchor hierarchy and its mapping to Go <code>Anchor</code> objects.", "hierarchy", "600px"),
	}
	return RenderHtml("Navigating and using the anchor hierarchy", Render(anchorBody, figs))
}

const anchorBody = `

<h2>Navigating and using the anchor hierarchy</h2>

<p>As you as obtain a <code>*Client</code> object <a href="api.html">after dialing</a>, you have your hands on
the root anchor. The type <code>*Client</code> implements the <code>Anchor</code> interface.

{{.FigHierarchy}}

<p>Anchors have two sets of methods: Those for navigating the anchor hierarchy and those for manipulating
elements attached to the anchor. We'll examine both types in turn.

<p>In general, it is safe to call all methods of anchor objects concurrently. In other words, you can
use the same anchor object from different goroutines without any synchronization.

<h3>Navigating</h3>

<p>Anchors have two methods pertaining to navigation of the anchor hierarchy:
<pre>
	View() map[string]Anchor
	Walk(walk []string) Anchor
</pre>

<p>Method <code>View</code> will return the current list of subanchors of this anchor. The result comes in the
form of a map from child names to anchor objects.

<p>Method <code>Walk</code> takes one <code>[]string</code> argument <code>walk</code>,
which is interpreted as a relative path (down the anchor hierarchy) from this anchor to a descendant. <code>Walk</code> traverses
the path and returns the descendant anchor, separated from this one by the path <code>walk</code>.

<p>Note that <code>Walk</code> always succeeds: If a child anchor is missing as <code>Walk</code> traverses down
the hierarchy, the anchor is created automatically. Anchors persist within the circuit cluster as long as they are
being used, otherwise they are garbage-collected. 
An anchor is in use if at least one of the following conditions is met:
<ul>
<li>It is being used by a client (in the form of a held <code>Anchor</code> object),
<li>It has an element attached to it (something like a process or a container, for example), or
<li>It has child anchors.
</ul>

<p>In other words, <code>Walk</code> is not only a method for accessing existing anchors but
also a method for creating new ones (should they not already exist). 

<p>The only exception to this rule is posed at the root anchor. The root anchor logically corresponds
to the cluster as a whole. The only allowed subanchors of the root anchor are the ones corresponding
to available circuit servers, and these anchors are created and removed automatically by the system.

<p>If <code>Walk</code> is invoked at the root anchor with a path argument, whose first 
element is not present in the hierarchy, the invokation will panic to indicate an error. This is not
a critical panic and one can safely recover from it and continue.

<h3>Manipulating elements</h3>

<p>The <code>Anchor</code> interface has a set of <code>Make…</code> methods,
each of which creates a new resource (process, container, etc.) and, if successful, atomically
attaches it to the anchor. (These methods would fail with a non-nil error, if the anchor
already has an element attached to it.) 
<pre>
	MakeChan(int) (Chan, error)
	MakeProc(Cmd) (Proc, error)
	MakeDocker(cdocker.Run) (cdocker.Container, error)
	MakeNameserver(string) (Nameserver, error)
	MakeOnJoin() (Subscription, error)
	MakeOnLeave() (Subscription, error)
</pre>
<p>The use of these methods is detailed in the following sections, dedicated to
<a href="api-process.html">processes</a>, 
<a href="api-container.html">containers</a>, 
<a href="api-channel.html">channels</a>, 
<a href="api-name.html">name servers</a> and
<a href="api-subscription.html">subscriptions</a>.

<p>Anchors have two generic methods for manipulating elements as well:
<pre>
	Get() interface{}
	Scrub()
</pre>
<p>The <code>Get</code> method will return the element currently associated with the
anchor. (This would be an object of type <code>Chan</code>, <code>Proc</code>,
	<code>Container</code>, <code>Nameserver</code>, <code>Server</code> or <code>Subscription</code>.)

<p>The <code>Scrub</code> method will terminate the operation of the element
attached to this anchor and will remove the element from the anchor.

<h3>Auxiliary methods</h3>

<p>Anchors have a couple of auxiliary methods to facilitate programming:
<pre>
	Addr() string
	Path() string
</pre>
<p>The method <code>Addr</code> returns the circuit address of the server that is hosting this anchor.
The returned value would be a string of the form <code>circuit://…</code>.

<p>The method <code>Path</code> will return the file-system notation of the anchor's path.
This would be a string looking like <code>"/X50faec8c2b5f6418/mysql/shard/1"</code>, for instance.

        `
