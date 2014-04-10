// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package rhunix

import (
	"strconv"
	"syscall"
	"time"
)

func readSysStat(syscall_Stat_t interface{}) (uid, gid string, atime time.Time, ver uint64) {
	st, ok := syscall_Stat_t.(*syscall.Stat_t)
	if !ok {
		return
	}
	uid, gid = strconv.FormatUint(uint64(st.Uid), 10), strconv.FormatUint(uint64(st.Gid), 10)
	atime = time.Unix(st.Atim.Unix())
	ver = 0 // No story for file version numbers in POSIX
	return
}
