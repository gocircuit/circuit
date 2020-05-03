// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package config facilitates parsing JSON data with $-style environment variable references
package config

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

func Parse(v interface{}, r io.Reader) error {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return ParseString(v, string(src))
}

func ParseString(v interface{}, src string) error {
	var err error
	q := []byte(src)
	q, err = rewrite(q)
	if err != nil {
		return err
	}
	return json.Unmarshal(q, v)
}

func ParseFile(v interface{}, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	src, err = rewrite(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(src, v)
}

func asvalue(s string) string {
	result, _ := json.Marshal(s)
	return string(result)
}

var fm = template.FuncMap{
	"env":    os.Getenv,
	"val":    asvalue,
	"repo":   printRepo,
	"os":     printGOOS,
	"goroot": printGOROOT,
}

func rewrite(src []byte) ([]byte, error) {
	t := template.New("").Funcs(fm)
	if _, err := t.Parse(string(src)); err != nil {
		return nil, err
	}
	var w bytes.Buffer
	if err := t.Execute(&w, nil); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// printGOOS returns the host OS keyword
func printGOOS() string {
	return runtime.GOOS
}

// printGOROOT returns `go env GOROOT` on the local host
func printGOROOT() string {
	return runtime.GOROOT()
}

// printRepo returns the closest parent directory that is a repo root, or the empty string.
func printRepo() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	r := parentRepos(wd)
	if len(r) == 0 {
		return ""
	}
	return r[0]
}

func parentRepos(wd string) []string {
	var r []string
	p := strings.Split(path.Clean(wd), "/")
	if len(p) == 0 || p[0] != "" {
		return nil
	}
	p = p[1:]
	for i := range p {
		g := path.Join(p[:len(p)-i]...)
		g = "/" + g
		if isRepo(g) {
			r = append(r, g)
		}
	}
	return r
}

func isRepo(p string) bool {
	if _, err := os.Stat(path.Join(p, ".git")); err == nil {
		return true
	}
	if _, err := os.Stat(path.Join(p, ".hg")); err == nil {
		return true
	}
	return false
}
