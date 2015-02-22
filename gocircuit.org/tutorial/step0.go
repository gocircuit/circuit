package tutorial

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderMysqlNodejsPage() string {
	return RenderHtml("Starting a MySQL and node.js stack using a circuit app", Render(installBody, nil))
}

const installBody = `
<h1>Starting a MySQL and node.js stack using a circuit app</h1>

<p>Start a fresh EC2 instance with an ubuntu base image. 

<h2>Prepare a host image</h2>

<pre>
	# apt-get update
	# apt-get install vim curl
</pre>

<h3>Install the circuit</h3>

<p>Start by installing the Go compiler and Git:

<pre>
	# apt-get install git golang
</pre>

<p>Create a temporary directory for building the circuit:

<pre>
	# mkdir -p /tmp/0/src && cd /tmp/0/src
	# declare -x GOPATH=/tmp/0
</pre>

<p>Fetch and build the circuit, then place the circuit executable in the system path:

<pre>
	# go get github.com/gocircuit/circuit/cmd/circuit
	# cp $GOPATH/bin/circuit /usr/local/bin
</pre>

<p>Next, configure the system to start the circuit daemon during the system booting sequence.

<p>??

<h3>Install MySQL server</h3>

<p>The installation will prompt you for a root user password — let's use <code>charlie</code>:

<pre>
	# apt-get install mysql-server
</pre>

<p>Start the server, so we can create a tutorial user and database:

<pre>
	# /etc/init.d/mysql start
</pre>

<p>Connect to the MySQL server as the administrator, using the password <code>charlie</code>:

<pre>
	# mysql -p
</pre>

<p>Create a user and a database, both named <code>tutorial</code>.

<pre>
	mysql> CREATE USER tutorial;
	mysql> CREATE DATABASE tutorial;
	mysql> GRANT ALL ON tutorial.*  TO tutorial;
</pre>

<p>Create table <code>Messages</code> for the tutorial application, after logging in as the <code>tutorial</code> user:

<pre>
	# mysql -u tutorial
	mysql> CREATE TABLE NameValue (name VARCHAR(100), value TEXT, PRIMARY KEY (name));
</pre>

<h3>Install node.js and the tutorial node.js app</h3>

<p>Install node.js:

<pre>
	# apt-get install nodejs
	# apt-get install npm
</pre>

<p>Install dependencies:

<pre>
	# npm install
</pre>


        `
