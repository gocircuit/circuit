FROM ubuntu:precise
MAINTAINER p@gocircuit.org

RUN echo "Building a docker image for a child circuit in a container"

# build.sh builds Golang's environment and installs the Circuit
ADD build.sh /go/util/build.sh
ADD start-circuit.sh /go/util/start-circuit.sh

RUN apt-get update -qq
RUN apt-get install -yqq mercurial git gcc

RUN /go/util/build.sh

CMD ["/go/util/start-circuit.sh"]
