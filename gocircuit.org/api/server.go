package api

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderServerPage() string {
	return RenderHtml("Using server", Render(serverBody, nil))
}

const serverBody = `

<h2>Using servers</h2>

<p>As we explain in the <a href="api.html">system abstraction</a> section, the
first level of anchors below the root (those that have paths like <code>/Xf1c8d96119cc6919</code>)
correspond to live hosts and we call them <em>server anchors</em>.

<p>Server anchors are automatically created and removed by the circuit system to
reflect the addition and removal (or death) of hosts.

<p>Every server anchors has a permanently attached <em>server element</em>. 
Server elements are objects whose methods allow you to retrieve some basic
information about the underlying host as well as the circuit daemon the manages the host.

<p>The server element for a given path can be accessed like any other element. For instance,
<pre>
	srv := root.Walk([]string{"Xf1c8d96119cc6919"})
</pre>

<p>Server elements implement the <code>Server</code> interface:
<pre>
	type Server interface {
		Profile(string) (io.ReadCloser, error)
		Peek() ServerStat
		Rejoin(string) error
		Suicide()
	}
</pre>

<h3>Retrieving server information</h3>

<p>The <code>Peek</code> method will return a structure containing the
circuit address of this server (which is a string of the form <code>circuit://…</code>),
and the time the server was launched according to the host's clock.
<pre>
	type ServerStat struct {
		Addr   string
		Joined time.Time
	}
</pre>

<p>Server elements cannot be scrubbed.

<h3>Merging circuit clusters</h3>

<p>The method <code>Rejoin</code> takes one argument which
is expected to be a target circuit address—i.e. a string of the form <code>circuit://…</code>.
Circuit addresses are a way of connecting to a circuit daemon directly.

<p><code>Rejoin</code> will cause the circuit daemon (underlying the server element)
to perform a join procedure with the target circuit daemon, specified by a circuit address.
This will result in merging this entire circuit cluster with the circuit cluster of the target.
If the target is already part of the same circuit cluster, no change will occur.

<h3>Miscellaneous</h3>

<p>The method <code>Suicide</code> will cause the circuit daemon, at the host
corresponding to the server element, to terminate itself.

<p>The method <code>Profile</code> will return profiling information for the circuit
daemon itself. Its argument specifies the name of the desired profile in the sense
of the Go package <code>runtime/pprof</code>.

        `
