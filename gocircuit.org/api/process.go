package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderProcessPage() string {
	figs := A{
		"FigMkProc": RenderFigurePngSvg("Process elements execute OS processes on behalf of the user.", "mkproc", "600px"),
	}
	return RenderHtml("Using processes", Render(processBody, figs))
}

const processBody = `

<h2>Using processes</h2>

<p>You can start an OS process on any host in your cluster by creating a
new <em>process element</em> at an anchor of your choosing that is a descendant of the
host's server anchor. The created process element becomes your interface to the
underlying OS process. 

<h3>Creating a process</h3>

<p>Suppose the variable <code>anchor</code> holds an <code>Anchor</code> object,
corresponding to a path in the anchor hierarchy that has no element attached to it.
For instance, say we obtained <code>anchor</code> like this:
<pre>
	anchor := root.Walk([]string{"Xe2ac4c8c83976ce6", "job", "demo"})
</pre>
<p>This anchor corresponds to the path <code>/Xe2ac4c8c83976ce6/job/demo</code>. 
(Read more on <a href="api-anchor.html">navigating anchors here</a>.)

<p>To create a new process element and attach it to <code>anchor</code>, 
we use the anchor's <code>MakeProc</code> method:
<pre>
	MakeProc(cmd Cmd) (Proc, error)
</pre>

<p><code>MakeProc</code> will start a new process on the host <code>/Xe2ac4c8c83976ce6</code>,
as specified by the command parameter <code>cmd</code>. If successful, it will create a 
corresponding process element and attach it to the anchor. <code>MakeProc</code> returns the 
newly created process element (of type <code>Proc</code>) as well as an 
<a href="api.html#errors">application error</a> (of type <code>error</code>), or it panics if a 
<a href="api.html#errors">system error</a> occurs.

<p>An application error can occur in one of two cases. Either the anchor already has another element
attached to it, or the process execution was rejected by the host OS (due to a missing binary or
	insufficient permissions, for example). 

<p><code>MakeProc</code> never blocks.

<p>The command parameter, of type <code>Cmd</code>, specifies the standard POSIX-level execution
parameters and an additional parameter called <code>Scrub</code>:
<pre>
type Cmd struct {
	Env []string
	Dir string
	Path string
	Args []string
	Scrub bool
}
</pre>

<p>If <code>Scrub</code> is set, the process element will automatically be detached from the anchor
and discarded, as soon as the underlying OS process exits. If <code>Scrub</code> is not set,
the process element will remain attached to the anchor even after the underlying OS process dies.
The latter regime is useful when one wants to start a job and return at a later time to check if
the job has already completed and what was its exit status. Furthermore, removing process elements
explicitly (rather than automatically) is a way of explicit accounting on the user's side. Thus
this regime is particularly well suited for applications that control circuit processes 
programmatically (as opposed to manually).


<h4>Example</h4>
<p>For instance, the following code executes the GNU list command:
<pre>
	proc, err := a.MakeProc(
		cli.Cmd{
			Env: []string{"TERM=xterm"},
			Dir: "/",
			Path: "/bin/ls",
			Args: []string{"-l", "/"},
			Scrub: true,
		},
	)
</pre>

<p>The following picture tries to illustrate the relationship between the
process element and the underlying OS process itself.

{{.FigMkProc}}

<h3>Controlling the standard file descriptors of a process</h3>

<p>After its invocation, <code>MakeProc</code> returns immediately,
while the underlying OS process is executing on the host machine.

<p>After a successful execution the user is obligated, by the POSIX 
standard, to take care of the standard input, output and error
streams of the underlying process. (For instance, if the standard
	input is not written to or closed, or if the output is not
	read from, some programs will pause in waiting.)

<p>The standard streams of the executed process can be retrieved
with the following methods of the process element:
<pre>
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
</pre>

<p>It is allowed to close the standard output and error at any point
into the stream. This will result in discarding all remaining data
in the stream, without blocking the underlying process.

<p>Eventually, the user is responsible for closing all standard streams
otherwise the underlying process will block and not exit.

<h3>Sending signals and killing processes</h3>

<p>You can send a POSIX signal to the underlying process
at any point (asynchronously) using:
<pre>
	Signal(sig string) error
</pre>

<p>The <code>sig</code> string must be one of the following recognized
signal names:
	<code>ABRT</code>,
	<code>ALRM</code>,
	<code>BUS</code>,
	<code>CHLD</code>,
	<code>CONT</code>,
	<code>FPE</code>,
	<code>HUP</code>,
	<code>ILL</code>,
	<code>INT</code>,
	<code>IO</code>,
	<code>IOT</code>,
	<code>KILL</code>,
	<code>PIPE</code>,
	<code>PROF</code>,
	<code>QUIT</code>,
	<code>SEGV</code>,
	<code>STOP</code>,
	<code>SYS</code>,
	<code>TERM</code>,
	<code>TRAP</code>,
	<code>TSTP</code>,
	<code>TTIN</code>,
	<code>TTOU</code>,
	<code>URG</code>,
	<code>USR1</code>,
	<code>USR2</code>,
	<code>VTALRM</code>,
	<code>WINCH</code>,
	<code>XCPU</code>,
	<code>XFSZ</code>.

<h3>Querying the status of a process asynchronously</h3>

<p>You can query the status of a process asynchronously, using:
<pre>
	Peek() ProcStat
</pre>

<p>The returned structure includes the command that started the process, a phase string describing the state of the
process and, in the event that the process has exited, an exit error value or <code>nil</code> on successful exit.
<pre>
	type ProcStat struct {
		Cmd Cmd
		Exit error
		Phase string
	}
</pre>

<p>The phase string takes on one of the following values:
<code>running</code>, 
<code>exited</code>,
<code>stopped</code>, 
<code>signaled</code>,
 <code>continued</code>.


<h3>Waiting until a process exits</h3>

<p>
<pre>
	Wait() (ProcStat, error)
</pre>

        `
