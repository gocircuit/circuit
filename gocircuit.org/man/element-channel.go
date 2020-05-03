package man

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderElementChannelPage() string {
	figs := A{
		"FigMkChan": RenderFigurePngSvg("Channel elements reside in the memory of a circuit server.", "mkchan", "600px"),
	}
	return RenderHtml("Circuit channel element", Render(elementChannelBody, figs))
}

const elementChannelBody = `

<h2>Example: Create a channel</h2>

<p>Again, take a look at what servers are available:

<pre>
	circuit ls /...
	/X88550014d4c82e4d
	/X938fe923bcdef2390
</pre>

<p>Pick one. Say <code>X88550014d4c82e4d</code>. Now, 
let's create a channel on <code>X88550014d4c82e4d</code>:

<pre>
	circuit mkchan /X88550014d4c82e4d/this/is/charlie 3
</pre>

<p>The last argument of this line is the channel buffer capacity,
analogously to the way channels are created in Go.

{{.FigMkChan}}

<p>Verify the channel was created:

<pre>
	circuit peek /X88550014d4c82e4d/this/is/charlie
</pre>

<p>This should print out something like this:

<pre>
	{
			"Cap": 3,
			"Closed": false,
			"Aborted": false,
			"NumSend": 0,
			"NumRecv": 0
	}
</pre>

<p>Sending a message to the channel is accomplished with the command

<pre>
	circuit send /X88550014d4c82e4d/this/is/charlie < some_file
</pre>

<p>The contents of the message is read out from the standard input of the
command above. This command will block until a receiver is available,
unless there is free space in the channel buffer for a message.

<p>When the command unblocks, it will send any data to the receiver.
If there is no receiver, but there is a space in the message buffer,
the command will also unblock and consume its standard input (saving
it for an eventual receiver) but only up to 32K bytes.

<p>Receiving is accomplished with the command

<pre>
	circuit recv /X88550014d4c82e4d/this/is/charlie
</pre>

<p>The received message will be produced on the standard output of 
the command above.

        `
