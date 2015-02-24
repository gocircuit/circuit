package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderProcessPage() string {
	figs := A{
		"FigMkProc": RenderFigurePngSvg("Process elements execute OS processes on behalf of the user.", "mkproc", "600px"),
	}
	return RenderHtml("Circuit API: Circuit process element", Render(processBody, figs))
}

const processBody = `

<h2>Using the Go API to control processes</h2>

<p>To create and manipulate process elements, one needs to import the circuit's Go client API:

<pre>
	import "github.com/gocircuit/circuit/client"
</pre>

<p>Let <code>cli</code>, a variable of type <code>*client.Client</code>, be an already established
connection to the circuit cluster. (<a href="api-client.html">How to connect a Go client to a circuit cluster.</a>)

<h3>Creating a process</h3>

<p>Suppose we already know that there are two servers in the circuit:

<pre>
	# circuit ls /
	/X88550014d4c82e4d
	/X938fe923bcdef2390
</pre>

<p>We would like to start a new process on the first server under the virtual path 
<code>/X88550014d4c82e4d/jobs/scrapy</code>. First, we need to obtain the
anchor for this virtual path:

<pre>
	a := cli.Walk([]string{"X88550014d4c82e4d", "jobs", "ls"})
</pre>

<p>The invocation of <code>Walk</code> always succeeds, as virtual paths are created
as needed (or otherwise they already exist as a circuit element is occupying them). The
invocation may only fail in panic, which is an indicator that this circuit server being
accessed, in this case <code>X88550014d4c82e4d</code>, has died or otherwise dropped
out of the cluster.

<p>Each anchor (virtual path) can have at most one circuit element (i.e. process, container, etc.)
attached to it. An anchor's <code>MakeProc</code> method will create a new
process element and attach it to the anchor:

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

        `
