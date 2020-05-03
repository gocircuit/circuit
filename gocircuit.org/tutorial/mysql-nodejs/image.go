package mysql_nodejs

import (
	. "github.com/hoijui/circuit/gocircuit.org/render"
)

func RenderImage() string {
	return RenderHtml("Prepare host images", Render(imageBody, nil))
}

const imageBody = `
<h1>Prepare host images</h1>

<p>We are going to describe here a sequence of steps that will
result in creating a new Amazon Machine Image (AMI), preloaded
with software needed for this tutorial.

<p>Start a fresh EC2 instance with an Ubuntu Linux base image. The plan
is to install the needed software in a few simple manual steps and then
save the state of the machine into the resulting image.

<h3>Install a few generic tools</h3>

<p>Begin by updating the packaging system and installing a few handy
generic tools:

<pre>
	# sudo apt-get update
	# sudo apt-get install vim curl git
</pre>

<h3>Install the circuit</h3>

<p>Installing the Go compiler:

<pre>
	# sudo apt-get install golang
</pre>

<p>Next, we need to create a directory for building the circuit.
This directory will serve as the <a href="https://golang.org/doc/code.html">GOPATH</a> directory
that points the Go compiler to a source tree.

<pre>
	# mkdir -p $HOME/0/src
	# echo "declare -x GOPATH=$HOME/0" >> ~/.bash_profile
	# source ~/.bash_profile
</pre>

<p>Fetch and build the circuit, then place the circuit executable in the system path:

<pre>
	# go get github.com/hoijui/circuit/cmd/circuit
	# cp $GOPATH/bin/circuit /usr/local/bin
</pre>

<p>Make sure the installation succeeded by running <code>circuit start</code>, to
start the circuit daemon, and then simply kill it with Control-C.

<p>Finally, prepare a scratch directory which we will later use for the circuit daemon
to write logging and other such information.

<pre>
	# mkdir /var/circuit
</pre>

<h3>Install MySQL server</h3>

<p>Install MySQL using the default packaged distribution.
The installation will prompt you for a root user password — 
feel free to use the empty string to simplify the tutorial:

<pre>
	# sudo apt-get install mysql-server
</pre>

<p>As a side-effect, the installer will put MySQL in the boot sequence of this machine.
We would like to disable that as we plan to manage (start/stop) the service through 
our circuit application. Thus, disable the automatic boot startup of MySQL using:

<pre>
	# echo manual | sudo tee /etc/init/mysql.override
</pre>

<h3>Install node.js and the tutorial node.js app</h3>

<p>Last, we need to install Node.js as well as the example Node.js app
that we have prepared for this tutorial, to serve as a RESTful HTTP API
front-end for the key/value store, backed by MySQL.

<p>Install Node.js with:

<pre>
	# sudo apt-get install nodejs npm
</pre>

<p>The Node.js app for this tutorial is located in the circuit's source tree.
We already have a local clone of the source tree (as a result of the <code>go get</code>
command earlier). We are simply going to copy it out of the source tree
and place it nearby for convenience:

<pre>
	# cd $HOME
	# cp -R $GOPATH/src/github.com/hoijui/circuit/tutorial/nodejs-using-mysql/nodejs-app .
</pre>

<p>And prepare the Node.js app for execution by fetching its Node.js package dependencies:

<pre>
	# cd $HOME
	# cd nodejs-app
	# npm install
</pre>

<h3>Save the image</h3>

<p>We are done installing. It is best to now “stop” (rather than “terminate”) the Amazon host instance.
Once the instance has stopped, use your Amazon EC2 console to save its current state to a new machine image.

        `
