package man

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderElementSubscriptionPage() string {
	return RenderHtml("Circuit subscription element", Render(subscriptionBody, nil))
}

const subscriptionBody = `

<h2>Example: Listen on server join and leave announcements</h2>

<p>The circuit provides two special element types <code>@join</code> and <code>@leave</code>, 
called <em>subscriptions</em>. Their job is to notify you when new
circuit servers join the systems or others leave it.
Both of them behave like receive-only channels.

<pre>
	circuit mk@join /X88550014d4c82e4d/watch/join
	circuit mk@leave /X88550014d4c82e4d/watch/leave
</pre>

<p>The join subscription delivers a new message each time a cicruit server joins
the system. The received message holds the anchor path of the new server.

<pre>
	circuit recv /X88550014d4c82e4d/watch/join
</pre>

<p>Similarly, the leave subscription delivers a new message each time
a circuit server disappears from the system.

<pre>
	circuit peek /X88550014d4c82e4d/watch/join
</pre>

        `
