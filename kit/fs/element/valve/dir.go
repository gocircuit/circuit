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

	A circuit channel is analogous to a “chan []byte” in Go.

INIT

	Initialize the channel by writing its desired buffer capacity
	to the "cap" file.

		echo 5 > cap

	This corresponds to “ch := make(chan []byte, 5)” in Go.

SEND

	To send a binary message through the channel, write the
	message to the "send" file.

		echo "¡hello, world!" >> send

	The open-file operation will block until the channel can
	accept the message.

	This is equivalent to “ch <- []byte("¡hello, world!")” in Go.

CLOSE

	Writing the text "close" to file "close" will close the channel permanently
	for writing. Buffered messages will still be receivable.

		echo "close" > close

RECV

	To receive the next channel message, read from the "recv" file.

		cat recv

	The open-file operation will block until a message is available.

	This is equivalent to “<-ch” in Go.

TRYING

	To send or receive a message without blocking, if possible,
	use the files "trysend" and "tryrecv" analogously to "send"
	and "recv".

	Opening a try-file will never block, but if the underlying
	channel operation is not available at the moment, the
	file-open operation will return with an error.

STAT

	Reading "stat" asynchronously retrieves channel runtime information.

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
