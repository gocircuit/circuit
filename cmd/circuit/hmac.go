// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/gocircuit/circuit/github.com/codegangsta/cli"
)

func keygen(c *cli.Context) string {
	rand.Seed(time.Now().UnixNano())
	seed := make([]byte, 32)
	for i, _ := range seed {
		seed[i] = byte(rand.Int31())
	}
	key := sha512.Sum512(seed)
	text := base64.StdEncoding.EncodeToString(key[:])
	
	return text
}

func keygenPrint(c *cli.Context) {

	fmt.Println(keygen(c))

}
