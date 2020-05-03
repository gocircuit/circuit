package man

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderInstallPage() string {
	return RenderHtml("Building and installing Circuit", Render(installBody, nil))
}

const installBody = `
<h1>Building and installing Circuit</h1>

<p>The Circuit comprises one small binary. It can be built for Linux and Darwin.

<p>Given that the <a href="http://golang.org">Go Language</a> compiler is <a href="http://golang.org/doc/install">installed</a>,
you can build and install the circuit binary with one line:

<pre>
	go get github.com/hoijui/circuit/cmd/circuit
</pre>

        `
