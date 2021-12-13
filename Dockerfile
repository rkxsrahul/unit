FROM golang:1.13


MAINTAINER Gursimran Singh <singhgursimran@me.com>

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin

ARG BUILD_ID
ENV BUILD_IMAGE=$BUILD_ID
ENV GO111MODULE=off


# build directories
ADD . /go/src/git.xenonstack.com/stacklabs/stacklabs-auth
WORKDIR /go/src/git.xenonstack.com/stacklabs/stacklabs-auth

# Go dep!
# RUN go get -u github.com/golang/dep/...
# RUN dep ensure -v

RUN go install git.xenonstack.com/stacklabs/stacklabs-auth
# ENTRYPOINT /go/bin/stacklabs-auth

EXPOSE 8000
