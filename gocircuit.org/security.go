package main

func RenderSecurityPage() string {
	return RenderHtml("Security and networking", Render(securityBody, nil))
}

const securityBody = `

<h2>Security</h2>

<p>By default, circuit servers and clients communicate over plaintext TCP.
A HMAC-based symmetric authentication, followed by an asymmetric
RC4 stream cipher is supported.

<p>To enable encryption, use the <code>-hmac</code> command-line option to point
the circuit executable to a file containing the private key for your circuit.
For instance:

<pre>
	circuit start -a 10.0.0.1 -hmac .hmac
</pre>

<p>Or, if you are invoking the tool:

<pre>
	circuit ls -hmac .hmac /...
</pre>

<p>Alternatively, you can set the environment <code>CIRCUIT_HMAC</code> to
point to the private key file.

<p>To generate a new private key for your circuit, use the command

<pre>
	circuit keygen
</pre>

<h2>Networking</h2>

<p>From a networking and protocol standpoint, circuit servers and
clients are peers: All communications (server-server and server-client)
use a common RPC framework which often entails a server
being able to reverse-dial into a client.

<p>For this reason, circuit clients (the circuit tool or your apps) CANNOT
be behind a firewall with respect to the servers they are dialing into.

        `
