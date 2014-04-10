// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package sel

import (
	"fmt"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// SelectDir 
type SelectDir struct {
	name string
	rh.FID
	dir *dir.Dir
	s   *Select
	rmv struct {
		sync.Mutex
		rmv func() // remove from parent
	}
}

func NewDir(name string, rmv func()) *SelectDir {
	d := &SelectDir{
		dir: dir.NewDir(),
		s:   MakeSelect(),
	}
	d.rmv.rmv = rmv
	d.FID = d.dir.FID()
	d.dir.AddChild("help", file.NewFileFID(file.NewByteReaderFile(
		func() []byte {
			return []byte(d.Help())
		}),
	))
	d.dir.AddChild("select", file.NewFileFID(NewSelectFile(d.s)))
	d.dir.AddChild("wait", file.NewFileFID(NewWaitFile(d.s)))
	d.dir.AddChild("error", file.NewFileFID(d.s.ErrorFile))
	return d
}

func (d *SelectDir) Walk(wname []string) (rh.FID, error) {
	if len(wname) > 0 {
		return d.FID.Walk(wname)
	}
	return d, nil
}

func (d *SelectDir) Remove() error {
	d.rmv.Lock()
	defer d.rmv.Unlock()
	if err := d.s.ClunkIfNotBusy(); err != nil {
		return rh.ErrBusy
	}
	if d.rmv.rmv != nil {
		d.rmv.rmv()
	}
	d.rmv.rmv = nil
	return nil
}

func (d *SelectDir) Help() string {
	return fmt.Sprintf(dirHelpFormat, d.name)
}

const dirHelpFormat = `
	This is the control directory for a select element named: %s

	Select elements are a mechanism for waiting on multiple
	files to open.

	All files in the circuit file system, which are named 
	"wait*" have the same behavior: They block on opening
	until some target event takes place. Select elements 
	were designed to wait for the first of a set of given
	wait-files to open.

SELECT

	To start a selection, write a JSON-encoded array of
	file names to the "select" file:

		echo << EOF > select
		[
			"/circuit/X130fc59d7291f8cf/dash/proc/paul/waitexit",
			"/circuit/X8d9ae9be389cfea7/dash/project1/chan/charlie/waitsend"
		]
		EOF

	File names should be absolute: They refer the local file system
	at the circuit worker, where the select is being created.

	Writing to "select" will return instantaneously, while 
	in the background the circuit will start waiting on opening
	all of the given files.

WAITING

	Once a selection has been started, trying to open the "wait"
	file will block until the first of the files being selected upon
	open successfully. The name of the latter file will then become
	readable as the contents of the "wait" file.

REUSE

	Select directories are reusable. After a selection has returned 
	from waiting, a new selection can be started, as in SELECT.

ERRORS

	Unsuccessful operations with the special in this directory will
	return file system errors. These errors are standardized in POSIX
	and are not descriptive for our purposes. To remedy this, after
	every file manipulation that returns a file system error, a detailed
	error message will be readable from the "error" file until the next
	file manipulation.

		cat error

`
