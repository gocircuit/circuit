package man

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderElementProcessPage() string {
	figs := A{
		"FigMkProc": RenderFigurePngSvg("Process elements execute OS processes on behalf of the user.", "mkproc", "600px"),
	}
	return RenderHtml("Circuit process element", Render(processBody, figs))
}

const processBody = `

<h2>Example: Make a process</h2>

<p>Here are a few examples. To run a new process on some chosen
cluster machine, first see what machines are available:

<pre>
	circuit ls /...
	/X88550014d4c82e4d
	/X938fe923bcdef2390
</pre>

<p>Run a new <code>ls</code> process:

<pre>
	circuit mkproc /X88550014d4c82e4d/pippi << EOF
	{
		"Path": "/bin/ls", 
		"Args":["/"]
	}
	EOF
</pre>

{{.FigMkProc}}

<p>See what happened:

<pre>
	circuit peek /X88550014d4c82e4d/pippi
</pre>

<p>Close the standard input to indicate no intention to write to it:

<pre>
	cat /dev/null | circuit stdin /X88550014d4c82e4d/pippi
</pre>

<p>Read the output (note that the output won't show until you close 
the standard input first, as shown above):

<pre>
	circuit stdout /X88550014d4c82e4d/pippi
</pre>

<p>Remove the process element from the anchor hierarchy

<pre>
	circuit scrub /X88550014d4c82e4d/pippi
</pre>

        `
