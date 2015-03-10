package main

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderIndexPage() string {
	figs := A{
		"FigFacade": RenderFigurePngSvg("Circuit API view into a cluster.", "facade", "700px"),
	}
	return RenderHtml(
		"Circuit: Self-managed infrastructure, programmatic monitoring and orchestration",
		Render(indexBody, figs),
	)
}

const indexBody = `

{{.FigFacade}}

<h1>Circuit: Self-managed cloud OS</h1>

<p>??

<h2>Sources</h2>

<p>Find the source repository for <a href="https://github.com/gocircuit/circuit">Circuit on GitHub</a>. 
Follow us on Twitter <a href="https://twitter.com/gocircuit">@gocircuit</a>.

<p>Submit <a href="">issues</a> to our GitHub repo. For discussions about using and developing
the Circuit visit <a href="https://groups.google.com/forum/#!forum/gocircuit-user">the Circuit User Group</a> and 
<a href="https://groups.google.com/forum/#!forum/gocircuit-dev">the Circuit Developer Group</a>, respectively.

<h2>Documentation</h2>

<ul>
<li><a href="install.html">Building and installing</a></li>
<li><a href="run.html">Running Circuit servers</a></li>
<li><a href="metaphor.html">Programming metaphor</a></li>
<li><a href="cmd.html">Command-line client</a>
<ul>
<li><a href="element-process.html">Process</a></li>
<li><a href="element-container.html">Container</a></li>
<li><a href="element-subscription.html">Subscription</a></li>
<li><a href="element-dns.html">Name server</a></li>
<li><a href="element-channel.html">Channel</a></li>
</ul>

<li><a href="api.html">Go client</a>
<ul>
<li><a href="api-anchor.html">Navigating the anchor hierarchy</a></li>
<li><a href="api-process.html">Process</a></li>
<li>Container</li>
<li>Subscription</li>
<li>Name server</li>
<li>Channel</li>
</ul>
</li>

</li>
<li><a href="security.html">Security and networking</a></li>
<li><a href="history.html">History and bibliography</a></li>
</ul>

<h2>Tutorials</h2>

<h3>Orchestrating a typical web app: Node.js using MySQL running on Amazon EC2</h3>
<ul>
<li><a href="tutorial-mysql-nodejs-overview.html">Overview</a></li>
<li><a href="tutorial-mysql-nodejs-image.html">Prepare host images</a></li>
<li><a href="tutorial-mysql-nodejs-boot.html">Boot the circuit cluster</a></li>
<li><a href="tutorial-mysql-nodejs-app.html">Write the circuit app</a></li>
<li><a href="tutorial-mysql-nodejs-run.html">Run the app on the cluster</a></li>
</ul>

<p>

        `
