// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

/*

	FORKING A GO ROUTINE ON A REMOTE RUNTIME

		import . "github.com/hoijui/circuit/pkg/use/circuit"

		type MyFunc struct{}
		func (MyFunc) AnyName(anyArg anyType) (anyReturn anyType) {
			...
		}
		func init() { types.RegisterFunc(MyFunc{}) }

		func main() {
			Go(conn, MyFunc{}, a1)
		}

*/
