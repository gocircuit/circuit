package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderChannelPage() string {
	return RenderHtml("Using server", Render(channelBody, nil))
}

const channelBody = `

<h2>Using channels</h2>

<p>Channel elements are FIFO queues for binary messages.
A channel physically lives on a particular host. It supports a ‘send’
operation which pushes a new binary message to the queue, and a 
‘receive’ operation which removes the next binary message from the
queue. 

<p>On creation, channels can be configured to buffer a fixed 
non-negative number of messages. When the buffer is full, send
operations block until a message is received. Receive operations,
on the other hand, block if there are no pending messages in the
buffer, until the next send operation is performed.

<p>Channels are created using the anchor's <code>MakeChan</code> method:
<pre>
	MakeChan(n int) (Chan, error)
</pre>

<p>An application error will be returned only if the anchor already has an
element attached to it. The integer parameter specifies the number of
messages the channel will buffer. This must be a non-negative integer.

<p>Interacting with the channel is done using the channel element's interface:
<pre>
	type Chan interface {
		Send() (io.WriteCloser, error)
		Recv() (io.ReadCloser, error)
		Close() error
		Stat() ChanStat
		Scrub()
	}
</pre>

<p>The <code>Scrub</code> method will forcefully abort and discard the
channel element, even if messages are pending in the channel's buffer.

<h3>Sending and receiving messages</h3>

<p>??

<h3>Channel closure</h3>

<p>??

<h3>Retrieving channel state</h3>
<pre>
	type ChanStat struct {
		Cap int
		Closed bool
		Aborted bool
		NumSend int
		NumRecv int
	}
</pre>

        `
