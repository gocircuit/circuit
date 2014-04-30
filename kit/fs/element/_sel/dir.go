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
	}
	d.s = MakeSelect(d)
	d.rmv.rmv = rmv
	d.FID = d.dir.FID()
	d.dir.AddChild("help", file.NewFileFID(file.NewByteReaderFile(
		func() []byte {
			return []byte(d.Help())
		}),
	))
	d.dir.AddChild("select", file.NewFileFID(NewSelectFile(d.s)))
	d.dir.AddChild("wait", file.NewFileFID(NewWaitFile(d.s)))
	d.dir.AddChild("trywait", file.NewFileFID(NewTryWaitFile(d.s)))
	d.dir.AddChild("abort", file.NewFileFID(NewAbortFile(d.s))) // TODO: Maybe aborting should be embodied in “rmdir”
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
	if err := d.s.Scrub(); err != nil {
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

	Select elements are a mechanism for waiting until one
	of multiple named files is ready to be opened.

	Select elements are inteded to be used in conjunction
	with other files in the circuit file system. (For instance,
	in a channel element directory, the "recv" file will block
	on opening until a message is available to be received.)

	When select unblocks, …

Initialization

	To start a selection, write a JSON-encoded array of file names
	to the "select" file:

		echo << EOF > select
		[
			{"op": "r", "file": "/circuit/X130fc59d7291f8cf/element/proc/paul/waitexit"},
			{"op": "w", "file": "/circuit/X8d9ae9be389cfea7/element/project1/chan/charlie/send"}
		]
		EOF

	File names should be absolute paths: They refer the local
	file system at the circuit worker, where the select element
	is being created.

	Writing to "select" will return instantaneously, while 
	in the background the circuit will start waiting on opening
	all of the given files. In the event of any error that
	prevents the selection from being performed, closing the 
	"select", after writing to it, will return an error.
	Human-readbale error messages can be read from the "error" file.

	(Note that some "echo" and/or shell implementations open a
	file twice, which would break the above example.)

Blocking Interface
	??

	Once a selection has been started, trying to open the "wait"
	file will block until the first of the files being selected upon
	opens successfully. Then the readable contents of "wait" will
	be a JSON structure, describing the unblocked file:

		cat wait
		{
			"clause": 0,
			"name":   "/circuit/X130fc59d7291f8cf/element/proc/paul/waitexit"
		}

	Where "clause" is the index of the file which unblocked on open,
	whereas "name" is its name.

Non-blocking Interface

	The file "trywait" is similar to "wait" in purpose, except that opening it
	is guaranteed not to block. If no selection event is ready, opening
	"trywait" will return an error.

Removal

	Select element directories can be removed with "rmdir" as long as 
	the none of the "select", "wait" or "trywait" files are currently open.

Errors

	Unsuccessful operations with the special in this directory will
	return file system errors. These errors are standardized in POSIX
	and are not descriptive for our purposes. To remedy this, after
	every file manipulation that returns a file system error, a detailed
	error message will be readable from the "error" file until the next
	file manipulation.

		cat error

`
