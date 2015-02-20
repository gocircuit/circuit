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
<li><a href="cmd.html">Command-line client</a></li>
<li>Elements
<ul>
<li><a href="element-process.html">Process</a></li>
<li><a href="element-container.html">Container</a></li>
<li><a href="element-subscription.html">Subscription</a></li>
<li><a href="element-dns.html">Name server</a></li>
<li><a href="element-channel.html">Channel</a></li>
</ul>
</li>
<li><a href="security.html">Security and networking</a></li>
<li><a href="history.html">History and bibliography</a></li>
</ul>

<h2>Tutorials</h2>

<h3>A typical web stack</h3>
<ul>
<li><a href="tutorial/nodejs-using-mysql/step0.html">Step 0: Prepare host VM images with the application software</a></li>
<li>Step 1: Starting a node.js and MySQL stack</li>
<li>Step 2: Starting a node.js, memcache and MySQL stack with co-location conditions</li>
<li>Step 3: Adding recovery logic for process failure</li>
<li>Step 4: Adding recovery from host failure</li>
</ul>

<!--h3>A server with a maintenance bot</h3>
<ul>
<li><a href="">Step 1: Starting a web server with a monitor-and-restart bot</a></li>
<li><a href="">Step 2: Controling the state of the bot through channels</a></li>
</ul-->

<p>

        `
