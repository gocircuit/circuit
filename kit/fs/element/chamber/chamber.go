// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chamber

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocircuit/circuit/kit/docker/docker"
	"github.com/gocircuit/circuit/use/n"
)

// The Chamber object manages a running child circuit worker, inside a Docker container.
type Chamber struct {
	ID    string
	Name  string
	Image string
	URL   string
	Addr  n.Addr
	sync.Mutex
	dc *docker.DockerCli
}

// The child container can use /var as private scratch space.
// The parent's binary directory is mounted under /vol/sup
//
//	PARENT					CHILD			DESCRIPTION
//	–––––––––––––––––––––––	–––––––––––––––	––––––––––––––––––––––––––––––––––––––
//	n/a						/var			Private child container scratch space
//	n/a						/circuit		Mount point for circuit file system
//	Binary directory		/genus/bin
//
func (t *Snowflake) NewChamber(name, image string) (chmb *Chamber, err error) {
	t.Lock()
	defer t.Unlock()
	//
	chmb = &Chamber{
		Name:  name,
		Image: image,
	}
	//
	var cmd = []string{
		"-d", "true", //  Start container in daemon mode
		"-name", name, // Set container name
		"-p", strconv.Itoa(t.Config.Port), // Expose container's circuit port to host
		"-v", "/var", // Create data volume for scratch space
		"-v", fmt.Sprintf("%s:%s", t.ParentDir, path.Join(t.Config.Genus, "bin")), // Share binary directory of this worker
		image, // docker image to use for this container
		//
		// circuit worker binary and arguments follow
		//
		path.Join(t.Config.Genus, "bin", t.ParentBinary), // Re-use the same binary that started this process, by accessing its mount point
		"-a", fmt.Sprintf(":%d", t.Config.Port), // Set worker port; the bind IP is (later computed) based on default route from `ip -s route`
		"-j", t.Config.ParentURL, // Join child worker to this worker
		"-m", "/circuit", // Mount circuit on worker local file system. XXX: Prepare FUSE within image
		"-xfs", "/", // Share all of local file system (XXX: for now)
		"-parent", "XXX", // Name of parent worker
	}
	// Run container
	if chmb.ID, _, err = t.dcRun(cmd...); err != nil {
		return nil, err
	}
	chmb.ID = strings.TrimSpace(chmb.ID)
	// Read child worker URL from its stdout
	if chmb.URL, _, err = t.dcLogs(chmb.ID); err != nil {
		defer t.dcKill(chmb.ID)
		defer t.dcScrub(chmb.ID)
		return nil, err
	}
	// Parse URL into circuit worker address and save it
	chmb.URL = strings.TrimSpace(chmb.URL)
	if chmb.Addr, err = n.ParseAddr(chmb.URL); err != nil {
		defer t.dcKill(chmb.ID)
		defer t.dcScrub(chmb.ID)
		return nil, err
	}

	// Retrieve external port for container by invoking `docker port CONTAINER PORT`
	// _, _, err := t.dcport(chmb.ID, strconv.Itoa(t.Config.Port))
	// if err != nil {
	// 	t.dcKill(chmb.ID)
	// 	t.dcScrub(chmb.ID)
	// 	return nil, err
	// }

	t.chamber[chmb.ID] = chmb
	go func() { // When container dies, scrub it form snowflake tables
		defer func() { // Scrubber
			t.Lock()
			defer t.Unlock()
			delete(t.chamber, chmb.ID)
		}()
		xx // wait for container to die
	}()
	//
	return chmb, nil
}

func (t *Snowflake) dc(f func(...string) error, arg ...string) (stdout, stderr string, err error) {
	t.buf.out.Reset()
	t.buf.err.Reset()
	if err = f(arg...); err != nil {
		return "", "", err
	}
	defer t.buf.out.Reset()
	defer t.buf.err.Reset()
	return t.buf.out.String(), t.buf.err.String(), nil
}

func (t *Snowflake) dcRun(arg ...string) (stdout, stderr string, err error) {
	return t.dc(t.dc.CmdRun, arg...)
}

func (t *Snowflake) dcPort(arg ...string) (stdout, stderr string, err error) {
	return t.dc(t.dc.CmdPort, arg...)
}

func (t *Snowflake) dcLogs(arg ...string) (stdout, stderr string, err error) {
	return t.dc(t.dc.CmdLogs, arg...)
}

func (t *Snowflake) dcKill(arg ...string) {
	t.dc(t.dc.CmdKill, arg...)
}

func (t *Snowflake) dcScrub(arg ...string) {
	t.dc(t.dc.CmdRm, arg...)
}

func (t *Snowflake) dcWait(arg ...string) (stdin, stdout string, err error) {
	t.dc(t.dc.CmdWait, arg...)
}
