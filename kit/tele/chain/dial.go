// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chain

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/gocircuit/circuit/kit/tele/trace"
)

type dialConn struct {
	dial dialFunc
	Conn
	x__     sync.Mutex
	ndial   SeqNo // Number of redials so far
	carrier net.Conn
}

type dialFunc func() (net.Conn, error)

func newDialConn(frame trace.Frame, id chainID, addr net.Addr, dial dialFunc, scrb func()) *dialConn {
	dc := &dialConn{dial: dial}
	dc.Conn.Start(frame, id, addr, (*dialLink)(dc), scrb)
	return dc
}

// dialLink is an alias for dialConn which implements the linker interface
type dialLink dialConn

// Link attempts to dial a new connection to the remote endpoint.
// If error is non-nil, the destination is permanently gone.
func (dl *dialLink) Link(reason error) (net.Conn, *bufio.Reader, SeqNo, error) {
	dl.x__.Lock()
	defer dl.x__.Unlock()

	if dl.carrier != nil {
		dl.carrier.Close()
	}
	if dl.carrier == nil && dl.ndial > 0 {
		return nil, nil, 0, io.ErrUnexpectedEOF
	}
	if dl.ndial > 0 {
		time.Sleep(CarrierRedialTimeout)
	}
	for {
		carrier, err := dl.dial()
		if err != nil {
			if err == ErrRIP {
				// Permanent error on the carrier connections dial attempts means
				// do not retry.
				dl.frame.Printf("permanently unreachable")
				return nil, nil, 0, err
			}
			// Non-permanent errors result in redial.
			dl.frame.Printf("redial because (%s)", err)
			time.Sleep(CarrierRedialTimeout + time.Duration((int64(CarrierRedialTimeout)*rand.Int63n(1e6))/1e6))
			continue
		}
		var broken bool
		if err, broken = dl.handshake(carrier); err != nil {
			if broken {
				log.Printf("dial rejected (%s)", err)
				return nil, nil, 0, ErrRIP
			}
			// All errors returned from dl.handshake where chainBroken is
			// false imply a broken carrier connection. We defer reporting the
			// broken connection error to the Conn object, which will spot it
			// when it tries to use the connection. This way, carrier-level
			// errors are treated uniformly within the Conn logic.
			// Furthermore, retrying from this error condition here would
			// break some test-only logic when the underlying sandbox
			// connection works in regime NOK=1 NDROP=0.
			log.Printf("carrier broke during handshake (%s)", err)
		}
		dl.carrier = carrier
		return dl.carrier, bufio.NewReader(carrier), dl.ndial, nil
	}
	panic("u")
}

func (dl *dialLink) handshake(carrier net.Conn) (err error, chainBroken bool) {
	var seqno SeqNo
	dl.ndial, seqno = dl.ndial+1, dl.ndial+1 // Connection count starts from 1
	msg := &msgDial{ID: dl.Conn.id, SeqNo: seqno}
	if err = msg.Write(carrier); err != nil {
		// Substrate error, chain is healthy
		return err, false
	}
	var welcome *msgWelcome
	if welcome, err = readMsgWelcome(carrier); err != nil {
		// Substrate error, chain is healthy
		return err, false
	}
	if welcome.Reject != rejectOK {
		// Reject from protocol, chain is permanently unacceptible from remote
		return errors.New(fmt.Sprintf("rejected, code %d", welcome.Reject)), true
	}
	return nil, false
}

// Kill shuts down the dialLink after a possibly concurring Link completes.
func (dl *dialLink) Kill() {
	dl.x__.Lock()
	defer dl.x__.Unlock()
	if dl.carrier == nil {
		return
	}
	dl.carrier.Close()
	dl.carrier = nil
}
