// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package prof

import (
	"bytes"
	"fmt"
)

type Stat struct {
	Type           string
	Begin          int64
	End            int64
	Abort          int64
	DurAvg         float64
	DurStdDev      float64
	AbortDurAvg    float64
	AbortDurStdDev float64
}

func (s *Stat) String() string {
	return fmt.Sprintf("type=%s completed=%10d pending=%10d avg/dev=%.1g/%.1g ms abort—avg/dev=%.1g/%.1g ms",
		s.Type, s.End, s.Begin-s.End-s.Abort,
		s.DurAvg/1e6, s.DurStdDev/1e6,
		s.AbortDurAvg/1e6, s.AbortDurStdDev/1e6,
	)
}

type WorkerStat struct {
	ReplyTotal *Stat
	ReplyProc  map[string]*Stat
	CallTotal  *Stat
	CallProc   map[string]*Stat
}

func (s *WorkerStat) String() string {
	var w bytes.Buffer
	for proc, stat := range s.ReplyProc {
		w.WriteString(proc)
		w.WriteString(": ")
		w.WriteString(stat.String())
		w.WriteByte('\n')
	}
	w.WriteString("———————————————————————————————————————————————————————————————————————————————\n")
	w.WriteString(s.ReplyTotal.String())

	for proc, stat := range s.CallProc {
		w.WriteString(proc)
		w.WriteString(": ")
		w.WriteString(stat.String())
		w.WriteByte('\n')
	}
	w.WriteString("———————————————————————————————————————————————————————————————————————————————\n")
	w.WriteString(s.CallTotal.String())
	return string(w.Bytes())
}
