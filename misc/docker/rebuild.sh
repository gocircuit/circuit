#!/bin/sh

mkdir -p /elektra
cd /elektra
hg clone https://code.google.com/p/go /elektra/go
cd /elektra/go/src
./all.bash

export PATH=$PATH:/elektra/go/bin:/elektra/0/bin:/elektra/aux/bin
export GOPATH=/elektra/0

mkdir -p /elektra/0/src

go get github.com/gocircuit/circuit/cmd/circuit
