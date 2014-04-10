// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package worker

import (
	"bufio"
	"bytes"
	"github.com/gocircuit/circuit/kit/iomisc"
	"github.com/gocircuit/circuit/load/config"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
)

// pie (Panic-If-Error) panics if err is non-nil
func pie(err interface{}) {
	if err != nil {
		panic(err)
	}
}

// pie2 panics if err is non-nil
func pie2(underscore interface{}, err interface{}) {
	pie(err)
}

// dbg is like a printf for debugging the interactions between
// daemonizer and runtime where stdandard out and error are not
// available to us to play with.
func dbg(n, s string) {
	cmd := exec.Command("sh")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic("huh")
	}
	cmd.Start()
	defer cmd.Wait()
	fmt.Fprintf(stdin, "echo '%s' >> /Users/petar/tmp/%s\n", s, n)
	stdin.Close()
}

func environLib() []string {
	var r []string
	for _, line := range os.Environ() {
		if strings.HasPrefix(line, "LD_") || strings.HasPrefix(line, "DYLD_") {
			r = append(r, line)
		}
	}
	return r
}

func daemonize() {
	// Change the file mode mask
	syscall.Umask(0)
	// Ensure that the process survives even if the invoking parent worker disconnects
	pie2(syscall.Setsid())
	// Change the current working directory.  This prevents the current
	// directory from being locked; hence not being able to remove it.
	pie(os.Chdir("/"))
}

func Daemonize(wc *config.WorkerConfig) {
	daemonize()

	// Make jail directory
	jail := path.Join(wc.Deploy.JailDir(), wc.Spark.ID.String())
	pie(os.MkdirAll(jail, 0700))

	// Prepare exec
	cmd := exec.Command(os.Args[0], wc.Spark.ID.String())
	cmd.Dir = jail
	cmd.Env = append(environLib(), fmt.Sprintf("%s=%s", config.RoleEnv, config.Worker))

	// Relay stdin of daemonizer to stdin of child runtime process
	var w bytes.Buffer
	pie(json.NewEncoder(&w).Encode(wc))
	cmd.Stdin = &w

	// Also save the config as a file for debugging purposes
	u, err := os.Create(path.Join(jail, "config"))
	if err != nil {
		panic(err)
	}
	pie(json.NewEncoder(u).Encode(wc))
	pie(u.Close())

	// Prepare out-of-band pipe for reading the child worker's PID and port
	oobReader, oobWriter, err := os.Pipe()
	pie(err)
	cmd.ExtraFiles = []*os.File{oobWriter}

	// Create stdout chain
	stdoutReader, err := cmd.StdoutPipe()
	pie(err)

	stdoutFile, err := os.Create(path.Join(jail, "out"))
	pie(err)
	defer stdoutFile.Close()

	go func() {
		io.Copy(iomisc.MultiWriter(stdoutFile, os.Stderr), iomisc.PrefixReader(fmt.Sprintf("%s:%s/out| ", wc.Spark.Host, wc.Spark.ID), stdoutReader))
	}()

	// Create stderr file
	stderrReader, err := cmd.StderrPipe()
	pie(err)

	stderrFile, err := os.Create(path.Join(jail, "err"))
	pie(err)
	defer stderrFile.Close()

	go func() {
		io.Copy(iomisc.MultiWriter(stderrFile, os.Stderr), iomisc.PrefixReader(fmt.Sprintf("%s:%s/err| ", wc.Spark.Host, wc.Spark.ID), stderrReader))
	}()

	// start
	pie(cmd.Start())

	// Read the first two lines of stdout. They should hold the Port and PID of the runtime process.
	back := bufio.NewReader(oobReader)

	// Read PID
	line, err := back.ReadString('\n')
	pie(err)

	pid, err := strconv.Atoi(strings.TrimSpace(line))
	pie(err)

	// Read port
	line, err = back.ReadString('\n')
	pie(err)

	port, err := strconv.Atoi(strings.TrimSpace(line))
	pie(err)

	// Close the pipe
	pie(oobReader.Close())

	if cmd.Process.Pid != pid {
		pie("pid mismatch")
	}

	fmt.Printf("%d\n%d\n", pid, port)
	// Sync is not supported on os.Stdout, at least on OSX
	// os.Stdout.Sync()

	// Survive kicker while the child worker is alive
	cmd.Wait()
	oobWriter.Close() // Close all pipes passed to the worker. Maybe unnecessary?

	os.Exit(0)
}
