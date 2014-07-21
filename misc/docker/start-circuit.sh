#!/bin/sh
# If another address is not specified through /go/util/addr
# we use the default one
if [ -f /go/util/addr ]; then
  /go/workspace/bin/circuit start -if eth0 -discover $(cat /go/util/addr)
else
  /go/workspace/bin/circuit start -if eth0 -discover 228.8.8.8:8788
fi
