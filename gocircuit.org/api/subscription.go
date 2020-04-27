package api

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderSubscriptionPage() string {
	return RenderHtml("Using subscriptions", Render(subscriptionBody, nil))
}

const subscriptionBody = `

<h2>Using subscriptions</h2>

<p>Subscriptions are a way of receiving notifications about events of a given type.
Presently, the circuit provides two types of subscriptions: 
<ul>
<li>Subscriptions for notifications about hosts joining the cluster, and
<li>Subscriptions for notifications about hosts leaving the cluster.
</ul>

<p>Just like process and other elements, a subscription is an element (a persistent object) that
is created and attached to an anchor. The methods of a subscription allow the user
to read events, one by one, in order of appearance, at the user's convenience.

<p>To create a new subscription element, use one of the following two <code>Anchor</code>
methods:
<pre>
	MakeOnJoin() (Subscription, error)
	MakeOnLeave() (Subscription, error)
</pre>

<p><code>MakeOnJoin</code> subscribes to the stream of ‘host joined the cluster’ events,
while <code>MakeOnLeave</code> subscribes to the stream of ‘host left the cluster’ events.
An application error will occur only if the underlying anchor is not free (i.e. it has an element
	already attached to it).

<p>The subscription element is represented by the following Go interface:

<pre>
type Subscription interface {
	Consume() (interface{}, bool)
	Peek() SubscriptionStat
	Scrub()
}
</pre>

<p>Subscriptions can be closed and discarded using the <code>Scrub</code> method of the
subscription element or of its anchor.

<h3>Consuming events</h3>

<p>Events are consumed using <code>Consume</code>. The first return value of <code>Consume</code>
holds the description of the event that was popped from the queue. 
The second return value is true if an event was successfully retrieved. Otherwise, the end-of-stream
has been reached permanently and the first return value will be nil.

<p>If the stream is still open and there are no events to be consumed, <code>Consume</code> will block.

<p>Host join and leave subscriptions return <code>string</code> events that
hold the textual path of the host that joined or left the network. These strings will
look like <code>/X36f63a7e4ae9df92</code>

<p>After a join-subscription is created, it will produce all the hosts that are currently 
in the cluster as events, and then it will continue producing new events as new hosts join later.

<p>After a leave-subscription is created, it will produce events only for hosts leaving the 
network after the subscription was created. Some leave events may be reported more than once.

<h3>Status of subscription queue</h3>

<p>The status of a subscription queue can be queried asynchronously using <code>Peek</code>.
The returned structure describes the type of the subscription, the number of pending (not yet consumed)
events, and whether the subscription has already been closed (by the user).
<pre>
type SubscriptionStat struct {
	Source string
	Pending int
	Closed bool
}
</pre>

<h3>Example</h3>

<p>Subscriptions, for join or leave events, are intended to be used via the following programming
pattern:
<pre>
join, err := MakeOnJoin()
if err != nil {
	…
}
for {
	event, ok := join.Consume()
	if !ok {
		…
	}
	host := event.(string)
	…
}
</pre>


        `
