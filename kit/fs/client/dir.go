// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"fmt"

	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/rh/xy"
	"github.com/gocircuit/circuit/kit/fs/server"
)

// ClientDir 
type ClientDir struct {
	home    *Peer        // peer structure of this worker
	peers   PeerFunc     //
	homeDir *server.ServerDir
	rh.FID             // namespace dir backing up this server dir
	dir *dir.Dir       // control object of namespace dir
}

// PeerFunc returns a list of known peers when called asynchronously.
type PeerFunc func() []*Peer

//
func NewDir(homeDir *server.ServerDir, home *Peer, peers PeerFunc) *ClientDir {
	d := &ClientDir{
		home:    home,
		peers:   peers,
		homeDir: homeDir,
		dir:     dir.NewDir(),
	}
	d.FID = d.dir.FID()
	d.dir.AddChild("help", newFile(d.Help))
	return d
}

func newFile(bodyFunc func() string) rh.FID {
	return file.NewFileFID(file.NewStringReaderFile(bodyFunc))
}

// Open synchronizes the list of known peers before opening the underlying directory structure.
func (d *ClientDir) Open(flag rh.Flag, intr rh.Intr) error {
	d.syncPeers()
	return d.FID.Open(flag, intr)
}

func (s *ClientDir) syncPeers() {
	s.dir.Clear()
	peers := s.peers()
	for _, peer := range peers {
		if peer.ID() == s.home.ID() {
			s.dir.AddChild(peer.Key(), s.homeDir) // will this play well with changing QIDs?
			continue
		}
		ysrv := xy.YServer{peer.Server}
		yssn, err := ysrv.SignIn("", "")
		if err != nil {
			continue
		}
		yfid, err := yssn.Walk(nil)
		if err != nil {
			continue
		}
		s.dir.AddChild(peer.Key(), yfid)
	}
}

func (s *ClientDir) Walk(wname []string) (rh.FID, error) {
	if len(wname) > 0 {
		return s.FID.Walk(wname)
	}
	return s, nil
}

func (d *ClientDir) Remove() error {
	return rh.ErrBusy
}

func (d *ClientDir) Help() string {
	kinstr := d.home.Kin.ID.String()
	return fmt.Sprintf(dirHelpFormat, kinstr)
}

const dirHelpFormat = `
	This directory ??: %s

	??

`
