package main

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderIndexPage() string {
	figs := A{
		"FigFacade": RenderFigurePngSvg(
			"Circuit API provides a dynamic third-person view of a compute cluster.", "facade", "700px"),
	}
	return RenderHtml(
		"Circuit: Self-managed infrastructure, programmatic monitoring and orchestration",
		Render(indexBody, figs),
	)
}

const indexBody = `

{{.FigFacade}}

<p>The circuit is a minimal distributed operating system that enables programmatic, reactive control
over hosts, processes and connections within a compute cluster.

<p>The circuit is unique in one respect: Once a circuit cluster is formed, the circuit system itself cannot 
fail—only individual hosts can. In contrast, comparable systems 
(like 
<a href="https://coreos.com/">CoreOS</a>, 
<a href="https://www.consul.io/">Consul</a> and 
<a href="http://mesosphere.com/">Mesosphere</a>)
can fail if the hardware hosting the system's own software fails.

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
<li><a href="api-connect.html">Connecting to a circuit cluster</a></li>
<li><a href="api-anchor.html">Navigating and using the anchor hierarchy</a></li>
<li><a href="api-process.html">Using processes</a></li>
<li>Using containers</li>
<li>Using subscription</li>
<li>Using name servers</li>
<li>Using channels</li>
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
