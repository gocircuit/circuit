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

<p>Notable here is the import of the circuit client package <code>"github.com/gocircuit/circuit/client"</code>
and the definition of function <code>fatalf()</code>, which we'll use to report terminal errors.


<h2>Connecting to the cluster</h2>

<p>The next step is to connect into a circuit server (specified by a circuit address) and obtain a client
object that will give us access to the circuit API.
The following subroutine takes a circuit address as an argument and returns a connected client object
as a result:

<pre>
func connect(addr string) *client.Client {
	defer func() {
		if r := recover(); r != nil {
			fatalf("could not connect: %v", r)
		}
	}()
	return client.Dial(addr, nil)
}
</pre>

<p>Note that by convention, the circuit client library universally reports loss 
of connection (or inability to establish connection) conditions via panics,
as they can occur anywhere in its methods.
Such panics are normal error conditions, and can be recovered from.
In our case, we prefer to terminate the program with an error message.

<p>The <code>main</code> function can now be updated to:
<pre>
func main() {
	flag.Parse()
	c := connect(*flagAddr)
	…
}
</pre>

<h2>Selecting hosts</h2>

<p>The next goal is to “list” the contents of the circuit cluster and to
choose two hosts out of the inventory — one for the MySQL database
and one for the Node.js front-end service.

<p>The circuit Go API represents all cluster resources in the form of one
big hierarchy of “anchors”. Each anchor can have any number of 
uniquely-named sub-anchors, and it can be associated with one
resource (process, server, container, etc.) Anchors are represented
by the interface type <code>client.Anchor</code>.

<p>The client object, of type <code>*client.Client</code>, is itself an
anchor (it implements <code>client.Anchor</code>) and it is in fact 
the root anchor of the circuit cluster's virtual hierarchy.

<p>The root anchor is unique in that it is not associated with any 
resource and its sub-anchors automatically exactly correspond
to the circuit servers that are presently members of the cluster.

<p>Every anchor has a <code>View</code> method:
<pre>
	View() map[string]client.Anchor
</pre>
which returns the sub-anchors of this anchor and their names.
If we invoke the <code>View</code> method of the root anchor,
we obtain a list of anchors corresponding to the currently live
circuit servers.

<p>We are going to use this to write a simple subroutine that blindly
picks a fixed number of hosts out of the available ones, re-using 
some hosts if necessary:

<pre>
func pickHosts(c *client.Client, n int) (hosts []client.Anchor) {
	defer func() {
		if recover() != nil {
			fatalf("client connection lost")
		}
	}()
	view := c.View()
	if len(view) == 0 {
		fatalf("no hosts in cluster")
	}
	for len(hosts) < n {
		for _, a := range view {
			if len(hosts) >= n {
				break
			}
			hosts = append(hosts, a)
		}
	}
	return
}
</pre>

<p>Note again here that a panic ensuing from <code>c.View()</code> would
indicate a broken connection between the client and the circuit server, in which
case we prefer to exit the program.

<p>We can further update <code>main</code> to:
<pre>
func main() {
	flag.Parse()
	c := connect(*flagAddr)
	hosts := pickHosts(c, 2)
	…
}
</pre>

<p>Before we continue with the main app logic—starting MySQL and starting Node.js—we
are going to make a small detour. We are going to implement a useful subroutine
that executes shell commands and scripts on any desired host directly from the
Go environment of our app.

<h2>A versatile orchestration subroutine</h2>

<p>The function <code>runShellStdin</code> takes an anchor parameter <code>host</code>,
which is expected to be an anchor corresponding to a circuit server. It
executes a desired shell command <code>cmd</code> on the corresponding host,
and it also supplies the string <code>stdin</code> as standard input to the
shell process.

<p>The function waits (blocks) until the shell process terminates and returns its
standard output in the form of a string. If the shell process exits in error, this
is reflected in a non-nil return error value. If the function fails due to loss of
connection (as opposed to due to an unsuccessful exit from the shell process),
<code>runShellStdin</code> will terminate the processes and exit with an error message.

<pre>
func runShellStdin(host client.Anchor, cmd, stdin string) (string, error) {
	defer func() {
		if recover() != nil {
			fatalf("connection to host lost")
		}
	}()
	job := host.Walk([]string{"shelljob", strconv.Itoa(rand.Int())})
	proc, _ := job.MakeProc(client.Cmd{
		Path:  "/bin/sh",
		Dir:   "/tmp",
		Args:  []string{"-c", cmd},
		Scrub: true,
	})
	go func() {
		io.Copy(proc.Stdin(), bytes.NewBufferString(stdin))
		proc.Stdin().Close() // Must close the standard input of the shell process.
	}()
	proc.Stderr().Close() // Close to indicate discarding standard error
	var buf bytes.Buffer
	io.Copy(&buf, proc.Stdout())
	stat, _ := proc.Wait()
	return buf.String(), stat.Exit
}
</pre>

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
