package mysql_nodejs

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderApp() string {
	return RenderHtml("Write the circuit app", Render(appBody, nil))
}

const appBody = `
<h1>Write the circuit app</h1>

<p>In this section we are going to write a Go program (the circuit app),
which when executed will connect to the circuit cluster and deploy
the tutorial's MySQL/Node.js-based web service using two hosts from
the cluster.

<p>The source of the finished app is located in:
<pre>
	$GOPATH/github.com/gocircuit/circuit/tutorial/nodejs-using-mysql/start-app/main.go
</pre>

<h2>First strokes</h2>

<p>This is going to be a single source-file Go program, which will expect exactly one command-line
parameter, <code>-addr</code>, specifying the circuit address to connect into. As
a start, the program will look like this:

<pre>
package main

import (
	"flag"
	"github.com/gocircuit/circuit/client"
)

var flagAddr = flag.String("addr", "", "circuit server address (looks like circuit://...)")

func fatalf(format string, arg ...interface{}) {
	println(fmt.Sprintf(format, arg...))
	os.Exit(1)
}

func main() {
	flag.Parse()
	…
}

</pre>

<h2>Connecting to the cluster</h2>

<p>

<p>

<h3>Start MySQL</h3>


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
	mysql> USE tutorial;
	mysql> CREATE TABLE NameValue (name VARCHAR(100), value TEXT, PRIMARY KEY (name));
</pre>


        `
