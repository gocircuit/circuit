package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderContainerPage() string {
	figs := A{
		"FigMkDkr": RenderFigurePngSvg("Docker elements are similar to processes.", "mkdkr", "600px"),
	}
	return RenderHtml("Using containers", Render(containerBody, figs))
}

const containerBody = `

<h2>Using containers</h2>

<p>The container-related types and structures of the circuit API are
in a dedicated package:
<pre>
	import "github.com/gocircuit/circuit/client/docker"
</pre>

<p>Container element manipulations and semantics are exactly analogous to 
their <a href="api-process.html">process element</a> counterparts.

<p>Given an anchor object, new containers are created using the anchor method:
<pre>
	MakeDocker(docker.Run) (docker.Container, error)
</pre>

<p>The <code>docker.Run</code> argument above is analogous to the <code>Cmd</code>
argument in the <code>MakeProc</code> method.
It specifies the execution parameters of the container, and is defined as:
<pre>
type Run struct {
	Image string
	Memory int64
	CpuShares int64
	Lxc []string
	Volume []string
	Dir string
	Entry string
	Env []string
	Path string
	Args []string
	Scrub bool
}
</pre>
<p>Excluding the field <code>Scrub</code>, all fields exactly match the standard 
Docker execution parameters which are explained in Docker's help:
<pre>
	docker help run
</pre>
<p>The field <code>Scrub</code> is also analogous to its counterpart in the process execution structure <code>Cmd</code>.
If <code>Scrub</code> is set, the container element will automatically be detached from the anchor and discarded, 
as soon as the underlying Docker container exits. 
If <code>Scrub</code> is not set, the container element will remain attached to the anchor even after the underlying 
Docker container dies.

<p>The methods of the container element interface are otherwise identical in form and meaning as those of the 
process element:
<pre>
type Container interface {
	Scrub()
	Peek() (*docker.Stat, error)
	Signal(sig string) error
	Wait() (*docker.Stat, error)
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
}
</pre>

<p>Finally, the <code>docker.Stat</code> structure (not shown here for space considerations) 
exactly captures all the container status variables that are available through the <code>docker inspect</code>
command.

<h4>Example</h4>
<p>The following snippet shows an example of creating a Docker container with an Ubuntu image,
which runs the <code>ls</code> command inside, while also specifying some resource limits and 
mapping some file system volumes:
<pre>
	proc, err := anchor.MakeDocker(
		docker.Run{
			Image: "ubuntu",
			Memory: 1000000000,
			CpuShares: 3,
			Volume: []string{"/webapp", "/src/webapp:/opt/webapp:ro"},
			Dir: "/",
			Path: "/bin/ls",
			Args: []string{"/"},
			Scrub: true,
		})
</pre>

{{.FigMkDkr}}

        `
