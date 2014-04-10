// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package fmt

import (
	"fmt"
)

func FormatBytes(n uint64) string {
	switch {
	case n < 1e3:
		return fmt.Sprintf("%dB", n)
	case n < 1e6:
		return fmt.Sprintf("%dKB", n/1e3)
	case n < 1e9:
		return fmt.Sprintf("%dMB", n/1e6)
	case n < 1e12:
		return fmt.Sprintf("%dGB", n/1e9)
	case n < 1e15:
		return fmt.Sprintf("%dTB", n/1e12)
	default:
		return fmt.Sprintf("%dPB", n/1e15)
	}
}
