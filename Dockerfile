FROM ubuntu:14.04
MAINTAINER Orne Brocaar <info@brocaar.com>

RUN apt-get update
RUN apt-get upgrade -y

# install docker
RUN apt-get install -y apt-transport-https
RUN sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
RUN sh -c "echo deb https://get.docker.io/ubuntu docker main > /etc/apt/sources.list.d/docker.list"
RUN apt-get update
RUN apt-get install -y lxc-docker
VOLUME /var/lib/docker

# install & setup golang
RUN apt-get install -y wget
RUN wget http://golang.org/dl/go1.3.linux-amd64.tar.gz
RUN tar -zxf go1.3.linux-amd64.tar.gz
RUN rm go1.3.linux-amd64.tar.gz
ENV PATH $PATH:/go/bin
ENV GOPATH /gocode
ENV GOROOT /go

# install dockerbuilder
RUN apt-get install -y make
RUN mkdir -p /gocode/src/github.com/brocaar
ADD . /gocode/src/github.com/brocaar/dockerbuilder
WORKDIR /gocode/src/github.com/brocaar/dockerbuilder
RUN make deps
RUN make
RUN mv dockerbuilder /usr/local/bin

EXPOSE 80
CMD /etc/init.d/docker start && /usr/local/bin/dockerbuilder
