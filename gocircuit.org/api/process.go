package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderProcessPage() string {
	figs := A{
		"FigMkProc": RenderFigurePngSvg("Process elements execute OS processes on behalf of the user.", "mkproc", "600px"),
	}
	return RenderHtml("Using processes", Render(processBody, figs))
}

const processBody = `

<h2>Using processes</h2>

<p>You can start an OS process on any host in your cluster by creating a
new <em>process element</em> at an anchor of your choosing that is a descendant of the
host's server anchor. The created process element becomes your interface to the
underlying OS process. 

<h3>Creating a process</h3>

<p>Suppose the variable <code>anchor</code> holds an <code>Anchor</code> object,
corresponding a path in the anchor hierarchy, which has no element attached to it.
For instance, you could have obtained the <code>anchor</code> from the root anchor
using something like this:
<pre>
	anchor := root.Walk([]string{"Xe2ac4c8c83976ce6", "job", "demo"})
</pre>
<p>This anchor corresponds to the path <code>/Xe2ac4c8c83976ce6/job/demo</code>. 
(See more on <a href="api-anchor.html">navigating anchors here</a>.)

<p>To create a new process element, use the anchor's <code>MakeProc</code> method:
<pre>
	MakeProc(cmd Cmd) (Proc, error)
</pre>

<p><code>MakeProc</code> will start a new process on the host <code>/Xe2ac4c8c83976ce6</code>,
as specified by the command parameter <code>cmd</code>. If successful, it will create a 
corresponding process element and attach it to the anchor. <code>MakeProc</code> returns the 
newly created process element (of type <code>Proc</code>) as well as an 
<a href="api.html#errors">application error</a>; or it panics if a 
<a href="api.html#errors">system error</a> occurs.

<p>

{{.FigMkProc}}

<h4>Example</h4>
<pre>
	proc, err := a.MakeProc(
		cli.Cmd{
			Env: []string{"TERM=xterm"},
			Dir: "/",
			Path: "/bin/ls",
			Args: []string{"-l", "/"},
			Scrub: true,
		},
	)
</pre>

<p>The returned error is non-nil if an element is already attached to the anchor <code>a</code> (i.e. to the path
<code>/X88550014d4c82e4d/jobs/ls</code> in our example).
Otherwise, 

<h3>Controlling the standard file descriptors of a process</h3>

<h3>Sending signals and killing processes</h3>

<h3>Querying the status of a process asynchronously</h3>

<h3>Waiting until a process exits</h3>

        `
