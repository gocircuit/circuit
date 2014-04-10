// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package worker implements spawning and killing of circuit workers on local and remote hosts
package worker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path"
	"strconv"

	"github.com/gocircuit/circuit/kit/posix"
	"github.com/gocircuit/circuit/load/config"
	"github.com/gocircuit/circuit/sys/tele"
	"github.com/gocircuit/circuit/use/errors"
	"github.com/gocircuit/circuit/use/n"
	"github.com/gocircuit/circuit/use/worker"
)

type Config struct {
	LibPath string
	Binary  string
	JailDir string
}

func New(libpath, binary, jaildir string) *Config {
	return &Config{
		LibPath: libpath,
		Binary:  binary,
		JailDir: jaildir,
	}
}

// (PARENT_HOST) --- parent worker --- ssh --- (CHILD_HOST) --- sh --- kicker --- child worker

func (c *Config) Spawn(h worker.Host, anchors ...string) (n.Addr, error) {

	host := h.(string)
	cmd := exec.Command("ssh", host, "sh")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	id := n.ChooseWorkerID()

	// Forward the stderr of the ssh process to this process' stderr
	posix.ForwardStderr(fmt.Sprintf("%s:kicker", id), stderr)

	// Start process
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	/*
		defer func() {
			go func() {
				// Use this spot to get notified when a child worker (or the connection to it) dies
				cmd.Wait()
			}()
		}()
	*/

	// Feed shell script to execute circuit binary
	bindir, _ := path.Split(c.Binary)
	if bindir == "" {
		panic("binary path not absolute")
	}
	var sh string
	if c.LibPath == "" {
		sh = fmt.Sprintf("cd %s\n%s=%s %s kicker:%s\n", bindir, config.RoleEnv, config.Daemonizer, c.Binary, id)
	} else {
		sh = fmt.Sprintf(
			"cd %s\nLD_LIBRARY_PATH=%s DYLD_LIBRARY_PATH=%s %s=%s %s\n",
			bindir, c.LibPath, c.LibPath, config.RoleEnv, config.Daemonizer, c.Binary)
	}
	stdin.Write([]byte(sh))

	// Read a chunk of magic from the child process, so as to know when to start feeding its STDIN.
	barrier := []byte(config.XXX_shell_barrier)
	if n, err := stdout.Read(barrier); n != len(config.XXX_shell_barrier) || err != nil || string(barrier) != config.XXX_shell_barrier {
		return nil, errors.NewError("kicker barrier not received")
	}

	// Write worker configuration to stdin of running worker process
	wc := &config.WorkerConfig{
		Spark: &config.SparkConfig{
			ID:       id,
			BindAddr: "0.0.0.0:0", // Empty string will mean don't listen at all
			Host:     host,
			Anchor:   append(anchors, fmt.Sprintf("/host/%s", host)),
		},
		Anchor: config.Config.Anchor,
		Deploy: config.Config.Deploy,
	}
	if err := json.NewEncoder(stdin).Encode(wc); err != nil {
		return nil, err
	}

	// Close stdin
	if err = stdin.Close(); err != nil {
		return nil, err
	}

	// Read the first two lines of stdout. They should hold the Port and PID of the runtime process.
	stdoutBuffer := bufio.NewReader(stdout)

	// First line equals PID
	line, err := stdoutBuffer.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = line[:len(line)-1]
	pid, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}

	// Second line equals port
	line, err = stdoutBuffer.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = line[:len(line)-1]
	port, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}

	netaddr, err := tele.ParseNetAddr(fmt.Sprintf("%s:%d", host, port)) // XXX: Tele-specific
	if err != nil {
		return nil, err
	}
	addr := tele.NewNetAddr(id, pid, netaddr)

	return addr, nil
}

func (c *Config) Kill(remote n.Addr) error {
	return kill(remote)
}

func kill(remote n.Addr) error {
	addr := remote.(*tele.Addr) // XXX: This is not clean. Assumes sys/tele. To be removed with upcoming new spawning mechanism.
	cmd := exec.Command("ssh", addr.TCP.IP.String(), "sh")

	stdinReader, stdinWriter := io.Pipe()
	cmd.Stdin = stdinReader

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(stdinWriter, "kill -KILL %d\n", addr.PID); err != nil {
		return err
	}
	stdinWriter.Close()

	return cmd.Wait()
}
