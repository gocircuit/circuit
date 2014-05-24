#!/usr/bin/env bash
apt-get update
apt-get install -y vim
apt-get install -y mercurial
apt-get install -y golang
mkdir -P $HOME/0/src # set GOPATH

# compile latest golang; apt-get golang is too old
hg clone https://code.google.com/p/go
cd go/src
./all.bash
