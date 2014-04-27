// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package server

import (
	"fmt"

	"github.com/gocircuit/circuit/kit/kinfolk"
	"github.com/gocircuit/circuit/kit/fs/element/dash"
	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// ServerDir 
type ServerDir struct {
	kin kinfolk.KinXID // kinfolk ID of the circuit owning this resource server
	rh.FID             // namespace dir backing up this server dir
	dir *dir.Dir       // control object of namespace dir
	fs   rh.FID        // local file system share, if available
}

type Resource struct {
	Name   string
	Source interface{} // rh.Server, rh.FID
}

func NewDir(kinXID kinfolk.KinXID, shared ...Resource) (*ServerDir, error) {
	d := &ServerDir{
		kin: kinXID,
		dir: dir.NewDir(),
	}
	kinstr := d.kin.ID.String()
	// Add builtin resources
	d.FID = d.dir.FID()
	d.dir.AddChild("help", newFile(d.Help))
	d.dir.AddChild("element", dash.NonRemovableDashDir{dash.NewDir(fmt.Sprintf("%s/element", kinstr))})
	d.dir.AddChild("sys", NewDebugDir())
	// Add shared resources
	for _, rsc := range shared {
		if err := d.addResource(rsc.Name, rsc.Source); err != nil {
			return nil, err
		}
	}
	//
	return d, nil
}

func newFile(bodyFunc func() string) rh.FID {
	return file.NewFileFID(file.NewStringReaderFile(bodyFunc))
}

func (s *ServerDir) addResource(name string, src interface{}) error {
	var fid rh.FID
	switch t := src.(type) {
	case rh.Server:
		ssn, err := t.SignIn("server", "")
		if err != nil {
			return err
		}
		fid, err = ssn.Walk(nil)
		if err != nil {
			return err
		}
	case rh.FID:
		fid = t
	default:
		return fmt.Errorf("not a resource")
	}
	if name == "fs" { // if the shared resource is the local file system
		s.fs = fid
	}
	_, err := s.dir.AddChild(name, fid)
	return err
}

func (s *ServerDir) Walk(wname []string) (rh.FID, error) {
	if len(wname) > 0 {
		return s.FID.Walk(wname)
	}
	return s, nil
}

func (d *ServerDir) Remove() error {
	return rh.ErrBusy
}

func (d *ServerDir) Help() string {
	kinstr := d.kin.ID.String()
	return fmt.Sprintf(dirHelpFormat, kinstr, kinstr)
}

const dirHelpFormat = `
	This directory encloses all resources of the circuit instance with ID: %s

DASH

	Go to the "dash" subdirectory to create and operate circuit elements,
	like channels, processes and selectors.

	All circuit elements created as descendants of "dash" will be hosted 
	on the circuit instance with ID %s.

		cd dash
		cat help

SYS

	If curious, the "sys" directory contains circuit worker runtime information.

`
