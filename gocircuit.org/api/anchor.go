package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
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


        `
