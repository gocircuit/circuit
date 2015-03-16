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

<p>Name servers maintain a set of unique <em>names</em>, together with a set of <em>records</em>
associated with each name.

<p>New records are added using the <code>Set</code> method. The argument of <code>Set</code> is
a DNS resource record in standard DNS notation. The syntax of these records are described in more detail
in the documentation of <a href="http://github.com/miekg/dns">github.com/miekg/dns</a>
as well as this <a href="http://miek.nl/posts/2014/Aug/16/go-dns-package/">related blog article</a>.

<p>Every DNS resource record, for instance <code>"miek.nl. 3600 IN MX 10 mx.miek.nl."</code>,
starts with the name that the record pertains to. Each invocation of the <code>Set</code> command
<em>adds</em> a record to the list of records pertaining to the respective name.
For instance, the following command adds the record <code>"miek.nl. 3600 IN MX 10 mx.miek.nl."</code>
to the name <code>"miek.nl."</code>
<pre>
	if err := ns.Set("miek.nl. 3600 IN MX 10 mx.miek.nl."); err != nil {
		… // DNS record cannot be recognized
	}
</pre>

<p>The <code>Unset</code> method removes <em>all</em> records associated with a given name. For instance,
<pre>
	ns.Unset("miek.nl.")
</pre>

<h3>Server status</h3>

<p>At any point, the user can asynchronously retrieve the current status of a name server element, 
using the <code>Peek</code> method of the element. The returned structure (shown below)
contains a textual representation of the DNS server's address, and a map of all names and their
associated lists of resource records.

<pre>
	type NameserverStat struct {
		Address string
		Records map[string][]string
	}
</pre>

        `
