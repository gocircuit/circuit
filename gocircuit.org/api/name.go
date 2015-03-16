package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderNamePage() string {
	return RenderHtml("Using name servers", Render(nameBody, nil))
}

const nameBody = `

<h2>Using name servers</h2>

<p>Name server elements are an easy way of creating and configuring 
lightweight DNS server dynamically.

<p>To create a name server element,
use the following method of the <code>Anchor</code> interface:
<pre>
	MakeNameserver(addr string) (Nameserver, error)
</pre>

<p>The creation of a new name server element results in starting
a lightweight DNS server (which is serviced by the circuit daemon itself)
on the respective host where the anchor lives.

<p>If the address argument is the empty string, the DNS server will pick
an available port to listen to. Otherwise, it will try to bind itself to <code>addr</code>.

<p>An application error may be returned if either (i) the underlying anchor is
already busy with another element, or (ii) the DNS server could not bind 
to the provided address parameter.

<p>Internally, the DNS server is implemented using <a href="http://github.com/miekg/dns">github.com/miekg/dns</a>.

<p>A name server can be stopped and discarded by either using the <code>Scrub</code> method of the
name server element itself, or by using the <code>Scrub</code> method of the anchor that the
name server element is attached to.

<p>Name server elements have a simple interface:
<pre>
	type Nameserver interface {
		Set(rr string) error
		Unset(name string)
		Peek() NameserverStat
		Scrub()
	}
</pre>

<h3>Manipulating records</h3>

<p>??
<pre>
	if err := ns.Set("miek.nl. 3600 IN MX 10 mx.miek.nl."); err != nil {
		… // DNS record cannot be recognized
	}
</pre>

<p>??
<pre>
	ns.Unset("miek.nl.")
</pre>

<h3>Server status</h3>

<p>??
<pre>
	type NameserverStat struct {
		Address string
		Records map[string][]string
	}
</pre>

        `
