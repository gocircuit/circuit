// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package rh

// Perm represents permissions and flags
//
// Permissions are kept in the low-order bits of the file mode: owner read/write/execute permission
// represented as 1 in bits 8, 7, and 6 respectively (using 0 to number the low order). The group
// permissions are in bits 5, 4, and 3, and the other permissions are in bits 2, 1, and 0.
type Perm uint16

// Permission bits base triplet
const (
	PermRead  Perm = 0x4
	PermWrite Perm = 0x2
	PermExec  Perm = 0x1
)

func (p Perm) String() string {
	return render(permchars, uint16(p))
}

var permchars = []bitchar{
	{0400, 'r'},
	{0, '-'},
	{0200, 'w'},
	{0, '-'},
	{0100, 'x'},
	{0, '-'},
	{0040, 'r'},
	{0, '-'},
	{0020, 'w'},
	{0, '-'},
	{0010, 'x'},
	{0, '-'},
	{0004, 'r'},
	{0, '-'},
	{0002, 'w'},
	{0, '-'},
	{0001, 'x'},
	{0, '-'},
}

type bitchar struct {
	bit uint16
	c   int
}

func render(chars []bitchar, p uint16) string {
	s := ""
	did := false
	for _, bc := range chars {
		if p&bc.bit != 0 {
			did = true
			s += string(bc.c)
		}
		if bc.bit == 0 {
			if !did {
				s += string(bc.c)
			}
			did = false
		}
	}
	return s
}
