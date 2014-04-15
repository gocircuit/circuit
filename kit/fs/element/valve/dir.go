// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"fmt"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/namespace/dir"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
	"github.com/gocircuit/circuit/kit/fs/rh"
)

// ValveDir 
type ValveDir struct {
	name string
	rh.FID
	dir *dir.Dir
	v   *Valve
	rmv struct {
		sync.Mutex
		rmv func() // remove from parent
	}
}

func NewDir(name string, rmv func()) *ValveDir {
	d := &ValveDir{
		dir: dir.NewDir(),
		v:   MakeValve(),
	}
	d.rmv.rmv = rmv
	d.FID = d.dir.FID()
	d.dir.AddChild("help", file.NewFileFID(file.NewByteReaderFile(
		func() []byte {
			return []byte(d.Help())
		}),
	))
	d.dir.AddChild("send", file.NewFileFID(NewSendFile(d.v)))
	d.dir.AddChild("recv", file.NewFileFID(NewRecvFile(d.v)))
	//
	d.dir.AddChild("trysend", file.NewFileFID(NewTrySendFile(d.v)))
	d.dir.AddChild("tryrecv", file.NewFileFID(NewTryRecvFile(d.v)))
	//
	d.dir.AddChild("waitsend", file.NewFileFID(NewWaitSendFile(d.v)))
	d.dir.AddChild("waitrecv", file.NewFileFID(NewWaitRecvFile(d.v)))
	//
	d.dir.AddChild("close", file.NewFileFID(NewCloseFile(d.v)))
	d.dir.AddChild("cap", file.NewFileFID(NewCapFile(d.v)))
	d.dir.AddChild("stat", file.NewFileFID(file.NewByteReaderFile(
		func() []byte {
			if stat := d.v.GetStat(); stat != nil {
				return []byte(stat.String())
			}
			return []byte("closed\n")
		}),
	))
	d.dir.AddChild("error", file.NewFileFID(d.v.ErrorFile))
	return d
}

func (s *ValveDir) Walk(wname []string) (rh.FID, error) {
	if len(wname) > 0 {
		return s.FID.Walk(wname)
	}
	return s, nil
}

func (d *ValveDir) Remove() error {
	d.rmv.Lock()
	defer d.rmv.Unlock()
	if d.rmv.rmv != nil {
		d.rmv.rmv()
	}
	d.rmv.rmv = nil
	return nil
}

func (d *ValveDir) Help() string {
	return fmt.Sprintf(dirHelpFormat, d.name)
}

const dirHelpFormat = `
	This is the control directory for a circuit channel named: %s

	A circuit channel matches concurring write-session to "send"
	with read-sessions from "recv", while no more than "cap"
	copies of "send" can be open at any one time.

SEND

	When the user tries to open "send" for writing, the file open
	operation will block until fewer than capacity-many copies of "send"
	are still open.

	When opening "send" succeeds, the "send" file session is still 
	not matched with a receiver: it is now in the "channel's buffer".
	The buffer accommodates at most capacity-many unmatched send sessions.

	Attempting to write to "send" will usually block write after open,
	until the send session is matched up with a reader of "recv".
	Then, the matched open copies of "send" and "receive" become an
	unbuffered pipe.

		echo "¡hello, world!" >> send

CLOSE

	Writing the text "close" to file "close" will close the channel permanently
	for writing. Pending send sessions will still be receivable.

		echo close > close

RECV

	When the user tries to open "recv" for reading, the file open
	operation will block until it can be matched to an available
	open copy of send in the channel buffer.

	When a match occurs, reading from "recv" unblocks.

		cat recv

	CAP

	The capacity of the channel buffer can be set dynamically by
	writing a non-negative integral capacity to "cap".

		echo 5 > cap

	Shrinking the buffer will not neglect pending/buffered send sessions.
	The current capacity can be retrieved by reading from "cap".

		cat cap

WAITING

	Trying to open "waitrecv" for reading will block until send sessions
	are available for receiving. 

		cat waitrecv

	Analogously, trying to open "waitsend" for reading will block until
	there is enough space in the channel buffer to accommodate the next
	"send" open operation immediately without blocking.

		cat waitsend

STAT

	Reading "stat" asynchronously retrieves channel dynamic statistics.

		cat stat

ERRORS

	Unsuccessful operations with the special in this directory will
	return file system errors. These errors are standardized in POSIX
	and are not descriptive for our purposes. To remedy this, after
	every file manipulation that returns a file system error, a detailed
	error message will be readable from the "error" file until the next
	file manipulation.

		cat error

`
