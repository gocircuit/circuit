// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package tele implements Teleport Transport which can overcome network outages without affecting endpoint logic.
package tele

/*

	The Teleport Transport networking stack:

	+----------------+
	|     BLEND      | Logical connection de/multiplexing over a single underlying connection.
	+----------------+
	|     CODEC      | Stateful encoding/decoding layer, e.g. gob/ProtoBuf/etc.
	+----------------+
	|    CARRIER     | Underlying transport, e.g. TCP/WebRTC/sandbox/etc.
	+----------------+

*/
