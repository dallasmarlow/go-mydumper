FROM ubuntu

RUN apt-get update && apt-get install -y build-essential \
    cmake \
    libglib2.0-dev \
    libmysqlclient-dev 

ADD https://launchpad.net/mydumper/0.6/0.6.2/+download/mydumper-0.6.2.tar.gz /usr/local/src
WORKDIR /usr/local/src/mydumper-0.6.2

RUN tar xfvz ../mydumper-0.6.2.tar.gz -C .. && \
    cmake . && \
    make && \
    make install && \
    make clean

WORKDIR /root
