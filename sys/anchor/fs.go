// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"bytes"
	"container/list"
	"path"
	"sort"
	"sync"

	use "github.com/gocircuit/circuit/use/anchorfs"
)

// System is a pure in-memory implementation of an Anchor File System.
// Writes commit synchronously, but can be called concurrently.
// Reads are concurrent.
type System struct {
	root  *Node // immutable field
	slash *Dir  // Directory handle to the root node
	sync.Mutex
	anchors map[string]*Anchor // WorkerID -> Anchor
}

// NewSystem returns an empty anchor file system
func NewSystem() *System {
	r := NewNode("", nil)
	return &System{
		root:    r,
		slash:   newDir(r.Open("")),
		anchors: make(map[string]*Anchor),
	}
}

func (fs *System) Remove(worker Payload) {
	// Writes are synchronous
	fs.Lock()
	defer fs.Unlock()

	// Find the anchor and remove it from the index
	name := worker.FileName()
	anchor := fs.anchors[name]
	if anchor == nil {
		return
	}
	delete(fs.anchors, name)

	// Lock all directory nodes so we can execute anchor removal atomically
	dnodes := make([]*Node, len(anchor.Files))
	for i, fn := range anchor.Files {
		dn := fn.Parent() // Never nil
		dn.StartTx()
		if !dn.Remove(fn.Name()) {
			panic("u")
		}
		dnodes[i] = dn
	}

	// Commit transactions and trim file system tree
	for _, dn := range dnodes {
		dn.CommitTx()
	}
	for _, dn := range dnodes {
		dn.Prune()
	}
}

type pathStops struct {
	Path  string
	Stops []string
}
type sortStops []pathStops

func (s sortStops) Len() int {
	return len(s)
}

func (s sortStops) Less(i, j int) bool {
	return s[i].Path < s[j].Path
}

func (s sortStops) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Create atomically creates anchors for the specified worker within the set of directories dirs.
func (fs *System) Create(dirs []string, pay Payload) (err error) {
	// Writes are synchronous
	fs.Lock()
	defer fs.Unlock()

	// Check for duplicate anchors
	wid := pay.FileName()
	if _, ok := fs.anchors[wid]; ok {
		return ErrDup
	}

	// Sanitize and split directories into parts
	stops := make(sortStops, 0, len(dirs))
	dup := make(map[string]struct{})
	for _, d := range dirs {
		stps, pth, err := use.SanitizeDir(d)
		if err != nil {
			return err
		}
		if _, ok := dup[pth]; ok {
			// Skip duplicate directories
			continue
		}
		dup[pth] = struct{}{}
		stops = append(stops, pathStops{Path: pth, Stops: stps})
	}
	sort.Sort(stops)

	// Find and lock all attachment points
	dnodes := make([]*Dir, len(stops))
	for i, st := range stops {
		dn := fs.mkdir(st.Stops)
		dnodes[i] = dn
	}
	for _, dn := range dnodes {
		dn.node().StartTx()
	}

	// Prepare file nodes for atomic insertion
	anchor := NewAnchor(pay)
	fnodes := make([]*Node, len(stops))
	for i, stp := range stops {
		fnodes[i] = NewNode(path.Join(stp.Path, wid), anchor)
	}
	anchor.Files = fnodes

	// Attach file nodes and close write transactions
	for i, dn := range dnodes {
		dn.node().Add(fnodes[i])
		// For read consistency, individual node commits must begin only after all nodes have started their transactions.
		dn.node().CommitTx()
	}
	for _, dn := range dnodes {
		dn.node().Prune()
	}

	// Index anchor
	fs.anchors[wid] = anchor

	return nil
}

// mkdir finds or creates the dirrectory node for the path comprised by parts.
func (fs *System) mkdir(parts []string) *Dir {
	p := fs.slash
	for i, d := range parts {
		switch q, err := p.OpenDir(d); err {
		case ErrExist:
			// Create directory if not existent
			p.node().StartTx()
			nn := NewNode("/"+path.Join(parts[:i+1]...), nil)
			q = newDir(nn.Open("")) // q becomes a *Dir handle for nn
			if !p.node().Add(nn) {
				panic("u")
			}
			p.node().CommitTx()
			p = q
		case nil:
			p = q
		default:
			panic("u")
		}
	}
	return p
}

func (fs *System) OpenFile(fullpath string) (*File, error) {
	dir, file := path.Split(fullpath)
	d, err := fs.OpenDir(dir)
	if err != nil {
		return nil, err
	}
	return d.OpenFile(file)
}

func (fs *System) OpenDir(fullpath string) (*Dir, error) {
	parts, _, err := use.SanitizeDir(fullpath)
	if err != nil {
		return nil, err
	}
	p := fs.slash
	for _, d := range parts {
		q, err := p.OpenDir(d)
		if err != nil {
			return nil, err
		}
		p = q
	}
	return p, nil
}

func (fs *System) Dump() string {
	var w bytes.Buffer
	var q list.List // Queue of directories
	q.PushBack("/")
	for q.Len() > 0 {
		e := q.Front()
		q.Remove(e)
		dpath := e.Value.(string)
		d, err := fs.OpenDir(dpath)
		if err != nil {
			// Dir doesn't exist any more
			continue
		}
		_, files, dirs := d.List()
		// Print filenames first
		for _, fpath := range files {
			w.WriteString(fpath)
			w.WriteRune('\n')
		}
		// Queue up directories
		for _, dpath := range dirs {
			q.PushBack(dpath)
		}
	}
	return string(w.Bytes())
}
