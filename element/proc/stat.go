// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"encoding/json"
	"strings"
	"syscall"
)

type Stat struct {
	Cmd   Cmd    `json:"cmd"`
	Exit  error  `json:"exit"`
	Phase string `json:"phase"`
}

func (s Stat) String() string {
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(0)
	}
	return string(b)
}

type Phase int

const (
	NotStarted Phase = iota
	Running
	Exited
	Stopped
	Signaled
	Continued
)

func (ph Phase) String() string {
	switch ph {
	case Running:
		return "running"
	case Exited:
		return "exited"
	case Stopped:
		return "stopped"
	case Signaled:
		return "signaled"
	case Continued:
		return "continued"
	}
	return "unknown"
}

func ParseSignal(sig string) (s syscall.Signal, ok bool) {
	s, ok = sigMap[strings.ToUpper(sig)]
	return
}
