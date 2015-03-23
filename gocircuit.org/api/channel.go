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

<p>For instance,
<pre>
	ch, err := anchor.MakeChan(0)
	if err != nil {
		… // this anchor is already busy
	}
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

<p>The send/receive mechanism transports <code>io.ReadWriteCloser</code>
pipes between the sender and the receiver. The sender gets the write end
and the receiver gets the read end of the pipe.

<p>To send a message, the user invokes <code>Send</code>. If the channel
has been closed or aborted, a non-nil application error is returned. Otherwise,
a new pipe is created. Send attempts to send the pipe itself through the
channel. It will block until there is space in the channel buffer or there is
a matching call to <code>Receive</code>.

<p>When <code>Send</code> unblocks it returns an <code>io.WriteCloser</code> wherein
the user writes the content of the message and closes the stream to indicate end
of message.
<pre>
	w, err := ch.Send()
	if err != nil {
		… // channel has been closed
	}
	w.Write(msg)
	w.Close()
</pre>

<p>To receive a message, the user invokes <code>Recv</code>. The operation blocks
until there is a pipe in the channel buffer or a matching call to <code>Send</code>. 
Receive will return an application error if the channel has been closed. Otherwise
it returns the reader end of the pipe received.

<h3>Channel closure</h3>

<p>Channels can be closed synchronously with the <code>Close</code> method.
The closure event will be communicated to the receiver after all pending messages
have been received. <code>Close</code> returns an application error only if the
channel has already been clsoed.

<p>An alternative method of closing the channel is to invoke <code>Scrub</code>.
This will abort the channel and discard any pending messages.

<h3>Retrieving channel state</h3>

<p>The <code>Stat</code> method can be called asynchronously to retrieve the
current state of the channel. The returned structure includes the channel capacity,
whether the channel has been closed or aborted, the number of messages sent
and received. (The difference of sent and received messages is the count of messages
	pending reception in the channel buffer.)
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
