#!/bin/sh
ip route get 8.8.8.8 | grep 8.8 | awk '{print $7}'
