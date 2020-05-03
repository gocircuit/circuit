// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

/*

					Ref			PermRef			exportRewrite		importRewrite
	-----------------------------------------------------------------------------
	T				*ref		*permref		=					=
	xptr			=			panic			ptrPtrMsg			n/a
	*ref			=			panic			ptrMsg				n/a
	xpermptr		=			=				permPtrPtrMsg		n/a
	*permref		panic		=				permPtrMsg			n/a
	-----------------------------------------------------------------------------
	*ptrMsg			n/a			n/a				n/a					*ptr
	*ptrPtrMsg		n/a			n/a				n/a					*ptr
	*permPtrMsg		n/a			n/a				n/a					*permptr
	*permPtrPtrMsg	n/a			n/a				n/a					*permptr


	USER VS RUNTIME TYPES

	X     ≈ xptr,		*_ref,		*_ptr,		*ptrMsg,		*ptrPtrMsg
	PermX ≈ xpermptr,	*_permref,	*_permptr,	*permPtrMsg,	*permPtrPtrMsg

*/
