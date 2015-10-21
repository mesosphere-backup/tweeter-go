FROM alpine:3.2
RUN apk add --update bash go bzr git mercurial subversion openssh-client ca-certificates build-base && \
    rm -rf /var/cache/apk/*

RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
WORKDIR /go