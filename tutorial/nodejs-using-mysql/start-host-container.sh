#!/bin/sh
docker run -P -i -v `pwd`:/root -t gocircuit:tutorial /bin/bash
# docker run -p 3306:3306 -i -v `pwd`:/root -t gocircuit:tutorial /bin/bash
