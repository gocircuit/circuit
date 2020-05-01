// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package hmac

type Addr string

func (a Addr) String() string {
	return string(a)
}

func (a Addr) Network() string {
	return "hmac/tcp"
}
