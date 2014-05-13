// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tube

/*
	Operating diagram:

	                              +––––+
	                              |USER|
	                              +––––+
	                                 |
	           (BulkRead/Write/BulkWrite/Lookup/Scrub/Forget)
	                                 |
	 ·+––––+                 +–––––+–V––+                 +–––––+·
	··|Tube|––(Bulk/Write)––>|XTube|Tube|––(Bulk/Write)––>|XTube|··
	 ·+––––+                 +–––––+––o–+                 +–––––+·
	                                  |
	                               +——o——+
	                               |TABLE| Tube stops propagation of ops, already reflected in the local table.
	                               +—————+

	 UPSTREAM –––––>–––––>–––––>–––––>–––––>–––––>–––––> DOWNSTREAM

*/
