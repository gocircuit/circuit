package mysql_nodejs

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderOverview() string {
	figs := A{
		"FigOverview": RenderFigurePngSvg(
			`Bubbles represent processes. The MySQL daemon and the left circuit are
			running on Host 1, for instance. MySQL and Node.js communicate via
			a TCP link and so do the two circuit processes. The circuit application 
			process can run anywhere: inside or outside of the cluster.`,
			"tutorial/mysql-nodejs",
			"600px",
		),
	}
	return RenderHtml("Starting a MySQL and node.js stack using a circuit app", Render(overviewBody, figs))
}

const overviewBody = `
<h1>Overview</h1>

<p>
In this tutorial we are going to build an example cloud-based key/value store, based MySQL and Node.js, 
using the circuit and Amazon EC2 as a host provider.

{{.FigOverview}}

<p>We are going to build and deploy the application step-by-step and from scratch. The high-level process
is as follows:

<ul>
<li>First, prepare an Ubuntu host image with all software that we might need 
(such as circuit, MySQL, Node.js, etc.).
<li>Second, launch multiple host instances based on the prepared image and
link them into a single circuit cluster.
<li>Third, implement a circuit app in the Go language, which orchestrates
the execution of the cloud-based key/value store, by assembling its components.
<li>Finally, launch the key/value store by connecting into the circuit cluster
and running the circuit app.
</ul>

<p>Steps 1 and 2 constitute the startup procedure for your circuit cluster.
Whereas steps 3 and 4 demonstrate the execution of a circuit application on the cluster.
Once the cluster is running, it can be reused for multiple executions and terminations
of circuit applications.

<p>To follow through this tutorial, you will need an <a href="http://aws.amazon.com/ec2/">Amazon EC2</a>
account to provision virtual host instances.

        `
