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

<p>Let's break down what this function accomplishes:

<ol>
<li><pre>
	defer func() {
		if recover() != nil {
			fatalf("connection to host lost")
		}
	}()
</pre>

<p>The <code>defer</code> statement catches panics that may arise from the circuit API calls.
By convention, any such panic indicates that either (i) the particular host we are manipulating 
(through the methods of the anchor object) has become disconnected from the cluster, or
(ii) our client has lost connection to the circuit server that it initially connected to, using <code>client.Dial</code>.

<p>In this example, we prefer to terminate the app if we encounter loss of connectivity of either kind.

<p>In general, one could detect whether (i) or (ii) was the cause for the panic.
For instance, if a subsequent call to a client method, like <code>View()</code>, also panics
then the client itself has been disconnected, i.e. condition (ii). In this case, you need to discard
the client object as well as any anchors derived from it. But if such a subsequent call does not panic, it implies
that the initial panic was caused by condition (i). In this case, only the host that your anchor
refers to has been disconnected and you can continue using the same client.

<li><pre>
	job := host.Walk([]string{"shelljob", strconv.Itoa(rand.Int())})
</pre>

<p>The next line, which invokes <code>host.Walk</code>, creates an anchor (i.e. a node in the 
	circuit's virtual hierarchy) for the shell process that we are about to execute.
	For instance, if the host anchor corresponds to a path like <code>/Xfea8b5b798f2fc09</code>,
	then the anchor <code>job</code> will correspond to a path like
	<code>/Xfea8b5b798f2fc09/shelljob/1234</code>, where <code>1234</code> is an
	integer that we pick randomly to make sure we arrive at an anchor that does not
	already have a resource attached to it.


<p>In general, calls to <code>anchor.Walk()</code> always succeed (as long as the implied
	underlying host is connected). If the anchor we are “walking” to does not already exist,
	it is automatically created. On the other hand, anchors that are not used by any clients
	and have no resources attached to them are eventually garbage-collected for you.

<li><pre>
	proc, _ := job.MakeProc(client.Cmd{
		Path:  "/bin/sh",
		Dir:   "/tmp",
		Args:  []string{"-c", cmd},
		Scrub: true,
	})
</pre>

<p>The following call to <code>job.MakeProc</code> executes the shell process
and creates a process handle — which we call a <em>process element</em> — and 
attaches the process element to the anchor <code>job</code>.
The process element is represented by the returned value in <code>proc</code>.
(In general, elements attached to an anchor can be retrieved using the <code>Get</code> method.)

<p>The function <code>MakeProc</code> returns as soon as the process is executed, it
does not wait for the process to complete. The returned error value, ignored in our example,
is non-nil only in the event that the binary to be executed is not found on the host.

<p>The argument to <code>MakeProc</code> specifies the command, as usual.
The field <code>Scrub</code>, when set to <code>true</code>, tells the circuit runtime to remove
the process anchor automatically when the process dies. (Normally anchors that
have resources attached to them are not garbage-collected from the virtual hierarchy.
They must be scrubbed explicitly by the user.)

<li><pre>
	go func() {
		io.Copy(proc.Stdin(), bytes.NewBufferString(stdin))
		proc.Stdin().Close() // Must close the standard input of the shell process.
	}()
	proc.Stderr().Close() // Close to indicate discarding standard error
	var buf bytes.Buffer
	io.Copy(&buf, proc.Stdout())
</pre>

<p>As soon as <code>MakeProc</code> returns, the process is running.
Our next goal is to take care of its standard streams: By POSIX convention,
every process will block if (i) it tries to read from standard input and there is nothing
to read and the descriptor is still open, or (ii) it tries to write to standard
error or output and they are not being consumed.

<p>We have direct access to the standard streams of the running process
via the methods
<pre>
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
</pre>
of the <code>proc</code> variable.

<p>In a separate goroutine, we write the contents of the parameter <code>stdin</code>
to the standard input of the shell process and then we close the stream, indicating that
no more input is to be expected.

<p>Meanwhile, in the main goroutine we first close the standard error stream.
This tells the circuit that all output to that stream should be discarded. Closing
a stream never blocks. 

<p>Finally, we block on greedily reading the standard output of the shell process
into a buffer until we encounter closure, i.e. an EOF condition. Closure of the
standard output stream happens immediately before the process exits. At
this point <code>io.Copy()</code> will unblock.

<li><pre>
	stat, _ := proc.Wait()
	return buf.String(), stat.Exit
</pre>

<p>At last we invoke <code>proc.Wait</code> to wait for the death of the
process and capture its exit state within the returned <code>stat</code> structure.
If the error field <code>stat.Exit</code> is non-nil, it means the process
exited in error.

</ol>

<p>Often we won't be interested in passing any data to the standard input of the shell process,
for which cases we add a shortcut subroutine:

<pre>
func runShell(host client.Anchor, cmd string) (string, error) {
	return runShellStdin(host, cmd, "")
}
</pre>

<p>We are now going to use <code>runShell</code> to fetch the public
and private IP addresses of any host on the cluster.

<h3>Retrieving EC2 host public and private IP addresses</h3>

<p>On any Amazon EC2 host instance, by definition, one is able to 
retrieve the public and private IP addresses of the host instance using
the following command-lines, respectively:

<pre>
curl http://169.254.169.254/latest/meta-data/public-ipv4
curl http://169.254.169.254/latest/meta-data/local-ipv4
</pre>

<p>Basing on that and using <code>runShell</code>, the following subroutine will
fetch the public address of any host on the circuit cluster, specified by its anchor:

<pre>
func getEc2PublicIP(host client.Anchor) string {
	out, err := runShell(host, "curl http://169.254.169.254/latest/meta-data/public-ipv4")
	if err != nil {
		fatalf("get ec2 public ip error: %v", err)
	}
	out = strings.TrimSpace(out)
	if _, err := net.ResolveIPAddr("ip", out); err != nil {
		fatalf("ip %q unrecognizable: %v", out, err)
	}
	return out
}
</pre>

<p>To retrieve the private host IP we implement a similar function <code>getEc2PrivateIP</code>,
which only differs from the above in that <code>public-ipv4</code> is substituted with 
<code>local-ipv4</code>.

<h2>Starting the MySQL database on host A</h2>

<p>We would like to write a routine that starts a fresh MySQL server on a given host and
returns its server address and port number as a result.

<pre>
	func startMysql(host client.Anchor) (ip, port string)
</pre>

<p>We are first going to describe the “manual” processs of starting a fresh
MySQL server, assuming we have a shell session at the host.

<p>Then we are going to show how this manual process can be codified
into a Go subroutine that performs its steps directly from the client application.

<h3>Manually starting MySQL at the host</h3>

<p>Let's asume you are at the shell of the host machine. The following
steps describe the way to start the MySQL server with a new database.

<p>Obtain the private IP address of this host:
<pre>
	$ IP=$(curl http://169.254.169.254/latest/meta-data/local-ipv4)
</pre>

<p>Rewrite MySQL's configuration file to bind to the private IP address
and the default port 3306:
<pre>
	$ sudo sed -i 's/^bind-address\s*=.*$/bind-address = '$IP:3306'/' /etc/mysql/my.cnf
</pre>

<p>Start the server:
<pre>
	$ sudo /etc/init.d/mysql start
</pre>

<p>Connect to MySQL as root to prepare the tutorial user and database:
<pre>
	$ sudo mysql
	mysql> DROP USER tutorial;
	mysql> DROP DATABASE tutorial;
	mysql> CREATE USER tutorial;
	mysql> CREATE DATABASE tutorial;
	mysql> GRANT ALL ON tutorial.*  TO tutorial;
</pre>

<p>Then connect as the <code>tutorial</code> user and set up the main table:
<pre>
	$ mysql -u tutorial
	mysql> USE tutorial;
	mysql> CREATE TABLE NameValue (name VARCHAR(100), value TEXT, PRIMARY KEY (name));
</pre>

<p>The database is now configured, up and accepting connections at <code>$IP:3306</code>.

<h3>Programmatically starting MySQL from the app</h3>

<p>Retracing the manual steps programmatically is straightforward, purely using the
subroutines <code>getEc2PrivateIP</code>, <code>runShell</code> and <code>runShellStdin</code>
that we created earlier.

<pre>
	func startMysql(host client.Anchor) (ip, port string) {

		// Retrieve the IP address of this host within the cluster's private network.
		ip = getEc2PrivateIP(host)

		// We use the default MySQL server port
		port = strconv.Itoa(3306)

		// Rewrite MySQL config to bind to the private host address
		cfg := fmt.Sprintf(` + "`" + `sudo sed -i 's/^bind-address\s*=.*$/bind-address = %s/' /etc/mysql/my.cnf` + "`" + `, ip)
		if _, err := runShell(host, cfg); err != nil {
			fatalf("mysql configuration error: %v", err)
		}

		// Start MySQL server
		if _, err := runShell(host, "sudo /etc/init.d/mysql start"); err != nil {
			fatalf("mysql start error: %v", err)
		}

		// Remove old database and user
		runShellStdin(host, "sudo /usr/bin/mysql", "DROP USER tutorial;")
		runShellStdin(host, "sudo /usr/bin/mysql", "DROP DATABASE tutorial;")

		// Create tutorial user and database within MySQL
		const m1 = ` + "`" + `
	CREATE USER tutorial;
	CREATE DATABASE tutorial;
	GRANT ALL ON tutorial.*  TO tutorial;
	` + "`" + `
		if _, err := runShellStdin(host, "sudo /usr/bin/mysql", m1); err != nil {
			fatalf("problem creating database and user: %v", err)
		}

		// Create key/value table within tutorial database
		const m2 = ` + "`" + `
	USE tutorial;
	CREATE TABLE NameValue (name VARCHAR(100), value TEXT, PRIMARY KEY (name));
	` + "`" + `
		if _, err := runShellStdin(host, "/usr/bin/mysql -u tutorial", m2); err != nil {
			fatalf("problem creating table: %v", err)
		}

		return
	}
</pre>

<p>We add a call to <code>startMysql</code> to the main logic:
<pre>
	func main() {
		flag.Parse()
		c := connect(*flagAddr)
		host := pickHosts(c, 2)

		mysqlIP, mysqlPort := startMysql(host[0])
		println("Started MySQL on private address:", mysqlIP, mysqlPort)
		…
	}
</pre>

<h2>Starting the Node.js app on host B</h2>

<p>Starting the Node.js app amounts to running the following command-line on the target host:
<pre>
	$ sudo /usr/bin/nodejs nodejs-app/index.js \
		--mysql_host $MYSQL_HOST --mysql_port $MYSQL_PORT \
		--api_host $API_HOST --api_port $API_PORT \
		&> /tmp/tutorial-nodejs.log
</pre>
<p>The app finds the backend MySQL server via the arguments <code>--mysql_host</code>
and <code>--mysql_port</code>. While it binds its HTTP API server to the address given by
the arguments <code>--api_host</code> and <code>--api_port</code>.

<p>The function <code>startNodejs</code> takes a target host parameter, 
as well as the host and port of the backing MySQL server. It starts the Node.js
app on the target host and returns the public IP address and port of the HTTP API endpoint.

<pre>
	func startNodejs(host client.Anchor, mysqlIP, mysqlPort string) (ip, port string) {
		defer func() {
			if recover() != nil {
				fatalf("connection to host lost")
			}
		}()

		// Start node.js application
		ip = getEc2PublicIP(host)
		port = "8080"
		job := host.Walk([]string{"nodejs"})
		shell := fmt.Sprintf(
			"sudo /usr/bin/nodejs index.js "+
				"--mysql_host %s --mysql_port %s --api_host %s --api_port %s "+
				"&> /tmp/tutorial-nodejs.log",
			mysqlIP, mysqlPort,
			"0.0.0.0", port,
		)
		proc, err := job.MakeProc(client.Cmd{
			Path:  "/bin/sh",
			Dir:   "/home/ubuntu/nodejs-app",
			Args:  []string{"-c", shell},
			Scrub: true,
		})
		if err != nil {
			fatalf("nodejs app already running")
		}
		proc.Stdin().Close()
		proc.Stdout().Close()
		proc.Stderr().Close()

		return
	}
</pre>

<p>Note how we run the server. We execute a shell process, which itself executes the Node.js app. 
The shell process, which is the one created by the circuit, is a long-running one. It will run
for as long as the child Node.js server is running.

<p>As soon as the process is executed <code>MakeProc</code> returns, but the process continues
executing in the background. We then close all of its standard streams as we don't intend them to be used.

<p>The process element is attached to anchor of the form <code>/X686ea8f7374e59a2/nodejs</code>.
This will allow you to find it in the future and check its state, for instance, by using the command-lines
<code>circuit ls</code> and <code>circuit peek</code>.

<p>At last, we tie this function into the main logic, which completes our circuit app:

<pre>
	func main() {
		flag.Parse()
		c := connect(*flagAddr)
		host := pickHosts(c, 2)

		mysqlIP, mysqlPort := startMysql(host[0])
		println("Started MySQL on private address:", mysqlIP, mysqlPort)

		nodejsIP, nodejsPort := startNodejs(host[1], mysqlIP, mysqlPort)
		println("Started Node.js service on public address:", nodejsIP, nodejsPort)
	}
</pre>

        `
