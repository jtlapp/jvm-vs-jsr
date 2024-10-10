FROM ubuntu:latest

ENV GO_VERSION=1.23.2
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
ENV BASE_APP_URL=http://api-service:8080

RUN apt-get update && \
    apt-get install -y \
        vim \
        curl \
        git \
        build-essential && \
    rm -rf /var/lib/apt/lists/*

RUN curl -OL https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz && \
    rm go$GO_VERSION.linux-amd64.tar.gz

RUN mkdir -p $GOPATH/bin

COPY root/ /root/
COPY src/ /src/

# remove builds made on my development machine
RUN find /src -type f -name "benchmark" -exec rm {} \;

WORKDIR /src

CMD ["/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait"]
