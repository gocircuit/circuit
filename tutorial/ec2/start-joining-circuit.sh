#!/bin/sh

# Save the EC2 private IP address of this host to a variable.
ip_address=` + "`" + `ifconfig eth0 | awk '/inet addr/ {split($2, a, ":"); print a[2] }'` + "`" + `

# Start the circuit server
/usr/local/bin/circuit start -a ${ip_address}:11022 -j $1 1> /var/circuit/address 2> /var/circuit/log &
