#!/bin/sh

mkdir -p /go
cd /go
hg clone https://code.google.com/p/go /go/golang
cd /go/golang/src
./all.bash

echo "export ADDRESS=228.8.8.8:8788" >> /etc/profile
echo "export GOPATH=/go/workspace" >> /etc/profile
echo "export PATH=$PATH:/go/golang/bin:/go/workspace/bin:/go/util" >> /etc/profile
export GOPATH=/go/workspace
export PATH=$PATH:/go/golang/bin:/go/workspace/bin:/go/util

mkdir -p /go/workspace/src

go get github.com/hoijui/circuit/cmd/circuit
