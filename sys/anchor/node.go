// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"errors"
	"path"
	"runtime"
	"sync"
	"time"

	csync "github.com/gocircuit/circuit/kit/sync"
	use "github.com/gocircuit/circuit/use/anchorfs"
)

var (
	ErrTimeout = errors.New("agfs timeout")
	ErrExist   = errors.New("agfs exist")
	ErrKind    = errors.New("agfs not file/dir")
	ErrDup     = errors.New("agfs dup dir")
)

// Payload stands for an object that can be stored within an anchor.
type Payload interface {
	FileName() string // Unique textual representation
}

// An Anchor represents a set of files that represent the same entity.
type Anchor struct {
	Created time.Time // Last time file was written
	Data    Payload
	Files   []*Node // File nodes refering to this anchor
}

func NewAnchor(data Payload) *Anchor {
	return &Anchor{
		Created: time.Now(),
		Data:    data,
	}
}

// Node is a single node in the file system hierarchy. It represents a file or a directory.
// If anchor is nil, the node is a directory; otherwise, a file refering to anchor.
// If the node is a file, it cannot have children. (Directories cannot be anchors.)
type Node struct {
	// Immutable
	Path   string
	Anchor *Anchor

	csync.WaitUntil // Change notifications

	sync.Mutex // When multiple nodes need locking, obtain locks in root-to-leaf order.
	rev        use.Rev
	parent     *Node
	created    time.Time
	updated    time.Time
	children   map[string]*Node
	nref       int
	scrubbed   bool
}

// name is the fully-qualified name of the node
func NewNode(abs string, anchor *Anchor) *Node {
	return &Node{
		Path:     abs,
		Anchor:   anchor,
		children: make(map[string]*Node),
	}
}

func (N *Node) link(parent *Node, stamp time.Time) {
	N.Lock()
	defer N.Unlock()
	N.parent, N.created = parent, stamp
}

func (N *Node) Name() string {
	_, name := path.Split(N.Path)
	return name
}

// Writing

// Child creations and deletions are performed within a transaction so all nodes,
// corresponding to an anchor can be updated at once.

// StartTx begins a new transaction
func (N *Node) StartTx() {
	N.Lock()
}

// name is the last component of node's fully-qualified name.
func (N *Node) Add(node *Node) bool {
	name := node.Name()
	m := N.children[name]
	if m != nil {
		return false
	}
	N.children[name] = node
	N.updated = time.Now()
	node.link(N, N.updated)
	N.rev++
	return true
}

// name is the last component of node's fully-qualified name.
func (N *Node) Remove(name string) bool {
	m := N.children[name]
	if m == nil {
		return false
	}
	m.scrubUnderLock()
	delete(N.children, name)
	N.updated = time.Now()
	N.rev++
	return true
}

// CommitTx commits a transaction
func (N *Node) CommitTx() {
	defer N.WaitUntil.Broadcast()
	N.Unlock()
}

// scrub removes the parent pointer from a Node, so it does not get used by mistake
func (N *Node) scrub() {
	N.Lock()
	defer N.Unlock()
	N.scrubUnderLock()
}

func (N *Node) scrubUnderLock() {
	if N.nref > 0 {
		panic("scrubbing a referenced node")
	}
	if N.scrubbed {
		panic("already scrubbed")
	}
	N.parent = nil
	N.updated = time.Now()
	N.rev++
}

// Prune locks the parent of this node and then scrubs the node, if eligible
func (N *Node) Prune() {
	N.Lock()
	P := N.parent
	N.Unlock()
	if P == nil || N.Anchor != nil {
		// Cannot prune the root or a file node
		return
	}
	// Lock parent, then child
	P.Lock()
	defer P.Unlock()
	N.Lock()
	defer N.Unlock()
	if len(N.children) > 0 || N.nref > 0 {
		return
	}
	P.Remove(N.Name())
	N.scrubUnderLock()
}

// Atomic reading

func (N *Node) Parent() *Node {
	N.Lock()
	defer N.Unlock()
	return N.parent
}

func (N *Node) IsFile() bool {
	N.Lock()
	defer N.Unlock()
	return N.Anchor != nil
}

func (N *Node) List() (files, dirs []*Handle, rev use.Rev) {
	N.Lock()
	defer N.Unlock()
	return N.listUnderLock()
}

func (N *Node) listUnderLock() (files, dirs []*Handle, rev use.Rev) {
	for _, n := range N.children {
		if n.IsFile() {
			files = append(files, newHandle(n))
		} else {
			dirs = append(dirs, newHandle(n))
		}
	}
	return files, dirs, N.rev
}

// Handles

// Rename open to
func (N *Node) Open(name string) *Handle {
	N.Lock()
	defer N.Unlock()
	if N.scrubbed {
		return nil
	}
	if name == "" {
		// Open this node itself
		N.nref++
		return newHandle(N)
	}
	n := N.children[name]
	if n == nil {
		return nil
	}
	return n.Open("")
}

// Handle refers to a Node ...
type Handle struct {
	*Node
	once sync.Once
}

func newHandle(n *Node) *Handle {
	if n == nil {
		return nil
	}
	h := &Handle{Node: n}
	runtime.SetFinalizer(h, func(x *Handle) {
		x.Close()
	})
	return h
}

func (h *Handle) Close() {
	h.once.Do(func() {
		h.Node.unref()
		h.Node = nil
	})
}

func (N *Node) unref() {
	defer N.Prune()
	N.Lock()
	defer N.Unlock()
	N.nref--
}

// Waiter/blocking methods

func (N *Node) Change(sinceRev use.Rev) (files, dirs []*Handle, rev use.Rev) {
	files, dirs, rev, _ = N.ChangeExpire(sinceRev, 0)
	return
}

func (N *Node) ChangeExpire(sinceRev use.Rev, expire time.Duration) (files, dirs []*Handle, rev use.Rev, err error) {
	var tmo <-chan time.Time
	if expire > 0 {
		tmo = time.NewTimer(expire).C
	}
	for {
		N.Lock()
		if N.scrubbed {
			panic("u")
		}
		if N.rev > sinceRev {
			defer N.Unlock()
			files, dirs, rev = N.listUnderLock()
			return files, dirs, rev, nil
		}
		w := N.WaitUntil.MakeWaiter()
		N.Unlock()
		//
		select {
		case <-w:
		case <-tmo:
			return nil, nil, 0, ErrTimeout
		}
	}
	panic("u")
}
