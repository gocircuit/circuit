// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package trace has the side effect of installing HTTP endpoints that report tracing information
package trace

import (
	"net/http"
	"runtime/pprof"
	"strconv"
)

func init() {
	http.HandleFunc("/_pprof", serveRuntimeProfile)
	http.HandleFunc("/_g", serveGoroutineProfile)
	http.HandleFunc("/_s", serveStackProfile)
}

func serveStackProfile(w http.ResponseWriter, r *http.Request) {
	prof := pprof.Lookup("goroutine")
	if prof == nil {
		http.Error(w, "unknown profile name", 400)
		return
	}
	prof.WriteTo(w, 2)
}

func serveGoroutineProfile(w http.ResponseWriter, r *http.Request) {
	prof := pprof.Lookup("goroutine")
	if prof == nil {
		http.Error(w, "unknown profile name", 400)
		return
	}
	prof.WriteTo(w, 1)
}

func serveRuntimeProfile(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("n")
	debug, err := strconv.Atoi(r.URL.Query().Get("d"))
	if err != nil {
		http.Error(w, "non-integer or missing debug flag", 400)
		return
	}

	prof := pprof.Lookup(name)
	if prof == nil {
		http.Error(w, "unknown profile name", 400)
		return
	}
	prof.WriteTo(w, debug)
}
