package man

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderElementContainerPage() string {
	figs := A{
		"FigMkDkr": RenderFigurePngSvg("Docker elements are similar to processes.", "mkdkr", "600px"),
	}
	return RenderHtml("Circuit container element", Render(containerBody, figs))
}

const containerBody = `

<h2>Example: Make a docker container</h2>

<p>Much like for the case of OS processes, the circuit can create, manage and synchronize 
<a href="http://www.docker.com">Docker</a> containers,
and attach the corresponding <em>docker elements</em> to a path in the anchor file system.

<p>To allow creation of docker elements, any individual server must be started 
with the <code>-docker</code> switch. For instance:

<pre>
	circuit start -if eth0 -discover 228.8.8.8:7711 -docker
</pre>

<p>To create and execute a new docker container, using the tool:

<pre>
	circuit mkdkr /X88550014d4c82e4d/docky << EOF
	{
		"Image": "ubuntu",
		"Memory": 1000000000,
		"CpuShares": 3,
		"Lxc": ["lxc.cgroup.cpuset.cpus = 0,1"],
		"Volume": ["/webapp", "/src/webapp:/opt/webapp:ro"],
		"Dir": "/",
		"Entry": "",
		"Env": ["PATH=/usr/bin"],
		"Path": "/bin/ls",
		"Args": ["/"],
	}
	EOF
</pre>

<p>Most of these fields can be omitted analogously to their command-line option counterparts 
of the <code>docker</code> command-line tool.

{{.FigMkDkr}}

<p>The remaining docker element commands are identical to those for processes:
<code>stdin</code>, <code>stdout</code>, <code>stderr</code>, <code>peek</code> and 
<code>wait</code>. In one exception, <code>peek</code> will return
a detailed description of the container, derived from <code>docker inspect</code>. 

        `
