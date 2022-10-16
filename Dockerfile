FROM ubuntu:latest

ARG GO_VERSION=1.19.2

RUN apt-get update -y 
RUN apt-get install -y wget curl make git libc-dev bash gcc
RUN wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"

RUN tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"
RUN rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR $GOPATH

COPY . $GOPATH/src

WORKDIR $GOPATH/src

RUN make install
RUN ln -sf /go/bin/OsmosisArbitrageBot /go/bin/defiant-swap

CMD [ "/go/bin/defiant-swap" ]
