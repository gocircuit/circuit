// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package chamber

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"os"
	"os/exec"
	"strings"
)

// GetLinkAddr returns the default external IP address of this machine.
//
// To compute it, GetLinkAddr picks out the first default IP address from
// `ip -s route`. Example outputs from `ip -s route`:
//
// default via 172.17.42.1 dev eth0
// 172.17.0.0/16 dev eth0  proto kernel  scope link  src 172.17.0.9
//
// default via 10.0.2.2 dev eth0
// default via 10.0.2.2 dev eth0  metric 100
// 10.0.0.0/24 dev eth1  proto kernel  scope link  src 10.0.0.10
// 10.0.0.0/24 dev eth2  proto kernel  scope link  src 10.0.0.10
// 10.0.0.0/24 dev eth3  proto kernel  scope link  src 10.0.0.10
// 10.0.2.0/24 dev eth0  proto kernel  scope link  src 10.0.2.15
// 10.0.3.0/24 dev lxcbr0  proto kernel  scope link  src 10.0.3.1
// 172.17.0.0/16 dev docker0  proto kernel  scope link  src 172.17.42.1
//
func GetLinkAddr() (ipaddr *net.IPAddr, err error) {
	out, err := exec.Command("ip", "-s", "route").CombinedOutput()
	if err != nil {
		return nil, err
	}
	var scanr = bufio.NewScanner(bytes.NewReader(out))
	for scanr.Scan() {
		words := strings.Split(scanr.Text(), " ")
		if len(words) == 0 || words[0] != "default" {
			continue
		}
		for _, w := range words[1:] {
			if len(strings.Split(w, ".")) != 4 {
				continue
			}
			ipaddr, err = net.ResolveIPAddr("ip", w)
			if err != nil {
				continue
			}
			return ipaddr, nil
		}
	}
	return nil, errors.New("cannot find default route")
}
