package tutorial

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderStep0Page() string {
	return RenderHtml("Step 0: Install app software on host images", Render(installBody, nil))
}

const installBody = `
<h1>Step 0: Install app software on host images</h1>

<p>Assuming an Ubuntu base image. Install Vim and node.js first.

<pre>
	# apt-get update
	# apt-get install vim
</pre>

<h2>Install MySQL server</h2>

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
	mysql> CREATE TABLE Messages (id INT NOT NULL AUTO_INCREMENT, at DATE, source VARCHAR(100), sink VARCHAR(100), body TEXT, PRIMARY KEY (id));
</pre>

<h2>Install node.js and a simple app that uses MySQL</h2>

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
