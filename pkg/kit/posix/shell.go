// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package posix

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/hoijui/circuit/pkg/kit/iomisc"
)

func ForwardCombined(prefix string, stderr, stdout io.Reader) {
	go func() {
		io.Copy(os.Stderr, iomisc.Combine(iomisc.PrefixReader(prefix+"/err| ", stderr), iomisc.PrefixReader(prefix+"/out| ", stdout)))
	}()
}

func ForwardStderr(prefix string, stderr io.Reader) {
	go func() {
		io.Copy(os.Stderr, iomisc.PrefixReader(prefix+"/err| ", stderr))
	}()
}

func Shell2(env Env, dir, shellScript string, show bool) error {
	cmd := exec.Command("sh", "-v")
	cmd.Dir = dir
	if env != nil {
		//println(fmt.Sprintf("%#v\n", env.Environ()))
		cmd.Env = env.Environ()
	}
	shellScript = "env | grep CGO 1>&2\nset -x\n" + shellScript
	PrintScript(dir, env, shellScript)
	cmd.Stdin = bytes.NewBufferString(shellScript)

	if show {
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		if err = cmd.Start(); err != nil {
			return err
		}
		// Build tool cannot write anything to stdout, other than the result directory at the end
		shell := "sh"
		io.Copy(os.Stderr, iomisc.Combine(iomisc.PrefixReader(shell+"/err| ", stderr), iomisc.PrefixReader(shell+"/out| ", stdout)))
	}
	return cmd.Wait()
}

func PrintScript(dir string, env Env, s string) {
	println("———————————————————————————————————————————", dir)

	for k, v := range env {
		println("% declare -x ", k+"="+v)
	}

	r := bytes.NewBufferString(s)
	for {
		line, err := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if line != "" {
			println("%", line)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err.Error())
		}
	}
	println("···········································")
}

// Env holds environment variables
type Env map[string]string

// Rename to GetEnv
func OSEnv() Env {
	environ := os.Environ()
	r := make(Env)
	for _, ev := range environ {
		kv := strings.SplitN(ev, "=", 2)
		if len(kv) != 2 {
			continue
		}
		r[kv[0]] = kv[1]
	}
	return r
}

func (env Env) Environ() []string {
	var r []string
	for k, v := range env {
		r = append(r, k+"="+v)
	}
	return r
}

func (env Env) Unset(key string) {
	delete(env, key)
}

func (env Env) Get(key string) string {
	return env[key]
}

func (env Env) Set(key, value string) {
	env[key] = value
}

func (env Env) Copy() Env {
	r := make(Env)
	for k, v := range env {
		r[k] = v
	}
	return r
}
