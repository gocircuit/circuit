package mysql_nodejs

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderMysqlNodejsApp() string {
	return RenderHtml("Starting a MySQL and node.js stack using a circuit app", Render(appBody, nil))
}

const appBody = `
<h1>Starting a MySQL and node.js stack using a circuit app</h1>

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
