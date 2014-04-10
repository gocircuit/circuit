// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

//
// #include <stdlib.h>
// #include <sys/param.h>
// #include <sys/mount.h>
//
import "C"

import (
	"unsafe"
)

type vfsconf struct {
	Name     string
	TypeNum  int
	RefCount int
}

func makeVFSConfig(vfc *C.struct_vfsconf) *vfsconf {
	return &vfsconf{
		Name:     C.GoStringN(&vfc.vfc_name[0], C.MFSNAMELEN),
		TypeNum:  int(vfc.vfc_typenum),
		RefCount: int(vfc.vfc_refcount),
	}
}

func getvfsbyname(name string) (*vfsconf, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var vfc C.struct_vfsconf
	if r, errno := C.getvfsbyname(cname, &vfc); r != 0 {
		return nil, errno
	}
	return makeVFSConfig(&vfc), nil
}
