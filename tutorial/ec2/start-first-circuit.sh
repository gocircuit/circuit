#!/bin/sh

# Save the EC2 private IP address of this host to a variable.
ip_address=`curl http://169.254.169.254/latest/meta-data/local-ipv4`

# Start the circuit server
/usr/local/bin/circuit start -a ${ip_address}:11022 1> /var/circuit/address 2> /var/circuit/log &
