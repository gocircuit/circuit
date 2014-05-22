// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/github.com/codegangsta/cli"
)

func fatalf(format string, arg ...interface{}) {
	println(fmt.Sprintf(format, arg...))
	os.Exit(1)
}

func readkey(x *cli.Context) (key []byte) {
	var hmac string
	switch {
	case x.IsSet("hmac"):
		hmac = x.String("hmac")
	case os.Getenv("CIRCUIT_HMAC") != "":
		hmac = os.Getenv("CIRCUIT_HMAC")
	default:
		return nil
	}
	b64, err := ioutil.ReadFile(hmac)
	if err != nil {
		fatalf("problem reading private key file (%s): %v", hmac, err)
	}
	if key, err = base64.StdEncoding.DecodeString(string(b64)); err != nil {
		fatalf("problem decoding base64 private key: %v", err)
	}
	return
}

func dial(x *cli.Context) *client.Client {
	var dialAddr string
	switch {
	case x.IsSet("dial"):
		dialAddr = x.String("dial")
	case os.Getenv("CIRCUIT") != "":
		buf, err := ioutil.ReadFile(os.Getenv("CIRCUIT"))
		if err != nil {
			fatalf("circuit environment file %s is not readable: %v", os.Getenv("CIRCUIT"), err)
		}
		dialAddr = strings.TrimSpace(string(buf))
	default:
		buf, err := ioutil.ReadFile(".circuit")
		if err != nil {
			fatalf("no dial address available; use flag -d or set CIRCUIT to a file name")
		}
		dialAddr = strings.TrimSpace(string(buf))
	}
	defer func() {
		if r := recover(); r != nil {
			fatalf("addressed server is gone or authentication failed")
		}
	}()
	return client.Dial(dialAddr, readkey(x))
}
