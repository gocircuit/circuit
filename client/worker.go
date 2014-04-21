// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"path"
)

// Worker is a client for a specific circuit worker's distributed control facilities.
type Worker struct {
	client *Client
	worker WorkerID
	dir    *Dir
	ns     *Namespace
}

// namespaceDir is the name of the root distributed programming directory
const namespaceDir = "namespace"

func newWorker(client *Client, worker WorkerID) (w *Worker, err error) {
	w = &Worker{
		client: client,
		worker: worker,
	}
	if w.dir, err = OpenDir(w.Path()); err != nil {
		return nil, err
	}
	if w.ns, err = newNamespace(w, nil); err != nil {
		return nil, err
	}
	return w, nil
}

// Path returns the path of this worker in the local circuit file system.
func (w *Worker) Path() string {
	return path.Join(w.client.Path(), string(w.worker))
}

// Namespace returns the given namespace directory of this worker.
func (w *Worker) Namespace(walk []string) (*Namespace, error) {
	return w.ns.Walk(walk)
}
