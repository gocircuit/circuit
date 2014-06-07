// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package kinfolk is an efficient “social” protocol for maintaining mutual awareness and sharing resources amongs circuit workers.
package kinfolk

/*

	The kin server maintains a “neighborhood” of peer kin servers.

	It reports peers that are discovered to be dead to the users.

	Kin hooks FolkXIDs with a trigger that catches Call panics

* Edit tube.go to use the folk's neighborhood methods

	Who is responsible for discovering dying neihbors and reporting to the tube?

		The FolkXIDs (of remote peers), that kin passes to its folk services,
		are rigged with a trigger that catches panics on Call and reports to the kin server.

		The Kin server reports discovered dead peers on its internal channel to Locus.

		Locus listens for RIP discoveries and writes them to the tube system.

*/
