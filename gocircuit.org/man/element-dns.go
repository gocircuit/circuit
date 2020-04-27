package man

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderElementDnsPage() string {
	return RenderHtml("Circuit DNS element", Render(dnsBody, nil))
}

const dnsBody = `

<h2>Example: Make a DNS server element</h2>

<p>Circuit allows you to create and dynamically configure one or more DNS
server elements on any circuit server.

<p>As before, pick an available circuit server, say <code>X88550014d4c82e4d</code>.
Create a new DNS server element, like so

<pre>
	circuit mkdns /X88550014d4c82e4d/mydns
</pre>

<p>This will start a new DNS server on the host of the given circuit server,
binding it to an available port. Alternatively, you can supply an
IP address argument specifying the bind address, as in

<pre>
	circuit mkdns /X88550014d4c82e4d/mydns 127.0.0.1:7711
</pre>

<p>Either way, you can always retrieve the address on which the DNS server
is listening by peeking into the corresponding circuit element:

<pre>
	circuit peek /X88550014d4c82e4d/mydns
</pre>

<p>This command will produce an output similar to this

<pre>
	{
	    "Address": "127.0.0.1:7711",
	    "Records": {}
	}
</pre>

<p>Once the DNS server element has been created, you can add resource records
to it, one at a time, using

<pre>
	circuit set /X88550014d4c82e4d/mydns "miek.nl. 3600 IN MX 10 mx.miek.nl."
</pre>

<p>Resource records use the canonical syntax, described in various RFCs.
You can find a list of such RFCs as well as examples in the DNS Go library
that underlies our implementation: <code>github.com/miekg/dns/dns.go</code>

<p>All records, associated with a given name can be removed with a single command:

<pre>
	circuit unset /X88550014d4c82e4d/mydns miek.nl.
</pre>

<p>The current set of active records can be retrieved by peeking into the element:

<pre>
	circuit peek /X88550014d4c82e4d/mydns
</pre>

<p>Assuming that a name has multiple records associated with it, peeking would produce
an output similar to this one:

<pre>
	{
		"Address": "127.0.0.1:7711",
		"Records": {
			"miek.nl.": [
				"miek.nl.\t3600\tIN\tMX\t10 mx.miek.nl.",
				"miek.nl.\t3600\tIN\tMX\t20 mx2.miek.nl."
			]
		}
	}
</pre>

        `
