FROM ubuntu:precise
MAINTAINER p@gocircuit.org

RUN echo "Building a docker image for a child circuit in a container"
ADD rebuild.sh /elektra/aux/bin/rebuild.sh
ADD getlinkaddr.sh /elektra/aux/bin/getlinkaddr.sh
ADD start-circuit.sh /elektra/aux/bin/start-circuit.sh

RUN apt-get update -qq
RUN apt-get install -yqq mercurial git gcc

RUN /elektra/aux/bin/rebuild.sh

# CMD circuit start -a ??:?? 
